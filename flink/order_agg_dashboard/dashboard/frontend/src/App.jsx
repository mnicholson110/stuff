import "maplibre-gl/dist/maplibre-gl.css";
import "./App.css";

import maplibregl from "maplibre-gl";
import React, { useEffect, useRef, useState } from "react";

const DARK_MAP_STYLE =
    "https://basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json";

export default function App() {
    const mapContainerRef = useRef(null);
    const mapRef = useRef(null);

    // Main store data (from snapshot and SSE)
    const [storeMap, setStoreMap] = useState({});

    // State to track animated stores.
    // Keys are store IDs and values are the animation start timestamp.
    const [orderAnimations, setOrderAnimations] = useState({});

    const [totalOrders, setTotalOrders] = useState(0);
    const [totalAmount, setTotalAmount] = useState(0);
    const [showOverlay, setShowOverlay] = useState(true);
    const [mapLoaded, setMapLoaded] = useState(false);

    // Initialize the map and add both sources and layers.
    useEffect(() => {
        if (mapRef.current) return;

        mapRef.current = new maplibregl.Map({
            container: mapContainerRef.current,
            style: DARK_MAP_STYLE,
            center: [-98.5795, 39.8283],
            zoom: 4,
        });

        mapRef.current.addControl(
            new maplibregl.NavigationControl(),
            "top-right",
        );

        mapRef.current.on("load", () => {
            // Add a source for the static store dots.
            mapRef.current.addSource("stores", {
                type: "geojson",
                data: {
                    type: "FeatureCollection",
                    features: [],
                },
            });

            // Add the static stores layer.
            mapRef.current.addLayer({
                id: "stores-layer",
                type: "circle",
                source: "stores",
                paint: {
                    "circle-radius": 5,
                    "circle-color": "blue",
                    "circle-stroke-width": 1,
                    "circle-stroke-color": "#fff",
                },
            });

            // --- New: Add source and layer for pulsing animations ---
            mapRef.current.addSource("storeAnimations", {
                type: "geojson",
                data: {
                    type: "FeatureCollection",
                    features: [],
                },
            });

            // Add the animations layer *behind* the static stores layer.
            mapRef.current.addLayer(
                {
                    id: "store-animations-layer",
                    type: "circle",
                    source: "storeAnimations",
                    paint: {
                        // Interpolate the radius from 5 to 15 based on progress.
                        "circle-radius": [
                            "interpolate",
                            ["linear"],
                            ["get", "progress"],
                            0,
                            5,
                            1,
                            15,
                        ],
                        // Fade the opacity from 0.8 down to 0.
                        "circle-opacity": [
                            "interpolate",
                            ["linear"],
                            ["get", "progress"],
                            0,
                            0.8,
                            1,
                            0,
                        ],
                        "circle-color": "red",
                        "circle-stroke-color": "#fff",
                        "circle-stroke-width": 1,
                    },
                },
                "stores-layer", // Specify the static layer as the reference so this layer is drawn beneath it.
            );
            // --------------------------------------------------------------

            setMapLoaded(true);

            // Fetch the initial snapshot.
            fetch("/snapshot")
                .then((res) => res.json())
                .then((data) => {
                    setStoreMap(data);
                })
                .catch((err) =>
                    console.error("Failed to fetch snapshot:", err),
                );
        });
    }, []);

    // SSE event handling: update storeMap and trigger animations.
    useEffect(() => {
        const es = new EventSource("/events");
        let updateQueue = [];

        es.onmessage = (event) => {
            if (!mapLoaded) return;
            try {
                const updatedStore = JSON.parse(event.data);
                updateQueue.push(updatedStore);
            } catch (err) {
                console.error("SSE parse error:", err, event.data);
            }
        };

        // Process one event at a time with a 50ms delay between events.
        const processQueue = () => {
            if (updateQueue.length > 0) {
                const store = updateQueue.shift();

                // Skip triggering an animation for tombstone records (from staleDataCheck)
                if (
                    !(store.order_count === 0 && store.total_order_amount === 0)
                ) {
                    setOrderAnimations((prev) => ({
                        ...prev,
                        [store.store_id]: Date.now(),
                    }));
                }
                // Update the main store data regardless.
                setStoreMap((prev) => ({
                    ...prev,
                    [store.store_id]: store,
                }));
            }
            // Schedule the next processing event (adjust delay as needed)
            setTimeout(processQueue, 50);
        };

        // Start processing the queue.
        processQueue();

        return () => {
            es.close();
        };
    }, [mapLoaded]);

    // Animation loop: update the animated layer based on the progress (0 to 1 over 500ms).
    useEffect(() => {
        let animationFrameId;

        function animate() {
            const now = Date.now();
            // Create a copy of the current animations.
            let updatedAnimations = { ...orderAnimations };
            const features = [];

            for (const [storeId, startTimestamp] of Object.entries(
                orderAnimations,
            )) {
                const elapsed = now - startTimestamp;
                const progress = Math.min(elapsed / 500, 1); // 500ms duration

                // If the animation has finished, remove this store from the animations.
                if (progress >= 1) {
                    delete updatedAnimations[storeId];
                } else {
                    const store = storeMap[storeId];
                    if (store && store.lat && store.lng) {
                        features.push({
                            type: "Feature",
                            geometry: {
                                type: "Point",
                                coordinates: [store.lng, store.lat],
                            },
                            properties: {
                                store_id: storeId,
                                progress, // used in the layer's paint expressions
                            },
                        });
                    }
                }
            }

            // Update the state if any animations have finished.
            if (
                Object.keys(updatedAnimations).length !==
                Object.keys(orderAnimations).length
            ) {
                setOrderAnimations(updatedAnimations);
            }

            // Update the animation layer's data.
            if (mapRef.current) {
                const animationSource =
                    mapRef.current.getSource("storeAnimations");
                if (animationSource) {
                    animationSource.setData({
                        type: "FeatureCollection",
                        features,
                    });
                }
            }
            animationFrameId = requestAnimationFrame(animate);
        }
        animationFrameId = requestAnimationFrame(animate);

        return () => cancelAnimationFrame(animationFrameId);
    }, [orderAnimations, storeMap]);

    // Update the static stores layer whenever storeMap changes.
    useEffect(() => {
        if (!mapRef.current) return;
        const src = mapRef.current.getSource("stores");
        if (!src) return;

        const features = Object.values(storeMap)
            .map((store) => {
                if (
                    !store.lat ||
                    !store.lng ||
                    !store.total_order_amount ||
                    !store.order_count
                )
                    return null;
                return {
                    type: "Feature",
                    geometry: {
                        type: "Point",
                        coordinates: [store.lng, store.lat],
                    },
                    properties: {
                        store_id: store.store_id,
                        order_count: store.order_count,
                        total_order_amount: store.total_order_amount,
                    },
                };
            })
            .filter(Boolean);

        src.setData({
            type: "FeatureCollection",
            features,
        });

        // Also update the aggregates overlay.
        fetch("/aggregates")
            .then((res) => res.json())
            .then((data) => {
                setTotalOrders(data.all_store_total_order_count);
                setTotalAmount(data.all_store_total_order_amount);
            })
            .catch((err) => console.error("Failed to fetch aggregates:", err));
    }, [storeMap]);

    // Popup on click for the static stores layer.
    useEffect(() => {
        if (!mapRef.current) return;
        mapRef.current.on("click", "stores-layer", (e) => {
            const feature = e.features[0];
            const { store_id, order_count, total_order_amount } =
                feature.properties;
            new maplibregl.Popup()
                .setLngLat(e.lngLat)
                .setHTML(
                    `
          <div style="color: #333;">
            <strong>Store ID:</strong> ${store_id}<br/>
            <strong>Orders:</strong> ${order_count}<br/>
            <strong>Amount:</strong> $${total_order_amount}
          </div>
        `,
                )
                .addTo(mapRef.current);
        });
    }, [mapLoaded]);

    return (
        <>
            <div className="map-container" ref={mapContainerRef} />

            {showOverlay && (
                <div className="data-overlay">
                    <h3>Aggregates</h3>
                    <p>
                        <strong>Total Orders: </strong>
                        {totalOrders.toLocaleString()}
                    </p>
                    <p>
                        <strong>Total Amount: </strong> $
                        {totalAmount.toLocaleString(undefined, {
                            minimumFractionDigits: 2,
                            maximumFractionDigits: 2,
                        })}
                    </p>
                    <button
                        className="toggle-btn"
                        onClick={() => setShowOverlay(false)}
                    >
                        Hide Data
                    </button>
                </div>
            )}

            {!showOverlay && (
                <div
                    className="data-overlay"
                    style={{ backgroundColor: "rgba(0,0,0,0.4)" }}
                >
                    <button
                        className="toggle-btn"
                        onClick={() => setShowOverlay(true)}
                    >
                        Show Data
                    </button>
                </div>
            )}
        </>
    );
}
