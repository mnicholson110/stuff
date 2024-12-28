import "maplibre-gl/dist/maplibre-gl.css";
import "./App.css";

import maplibregl from "maplibre-gl";
import React, { useEffect, useRef, useState } from "react";

const DARK_MAP_STYLE =
    "https://basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json";

function App() {
    const mapContainerRef = useRef(null);
    const mapRef = useRef(null);

    const [mapLoaded, setMapLoaded] = useState(true);
    const [storeMap, setStoreMap] = useState({});
    const [totalOrders, setTotalOrders] = useState(0);
    const [totalAmount, setTotalAmount] = useState(0);
    const [showOverlay, setShowOverlay] = useState(true);

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
            mapRef.current.addSource("stores", {
                type: "geojson",
                data: {
                    type: "FeatureCollection",
                    features: [],
                },
            });

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

            // popup on click
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

            setMapLoaded(true);
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

        // every 200ms, flush the queue in one go
        const flushInterval = setInterval(() => {
            if (updateQueue.length > 0) {
                setStoreMap((prevMap) => {
                    const newMap = { ...prevMap };

                    for (const store of updateQueue) {
                        newMap[store.store_id] = store;
                    }
                    updateQueue = [];

                    return newMap;
                });
            }
        }, 200);

        return () => {
            clearInterval(flushInterval);
            es.close();
        };
    }, []);

    useEffect(() => {
        if (!mapRef.current) return;

        const src = mapRef.current.getSource("stores");
        if (!src) return;

        const features = Object.values(storeMap)
            .map((store) => {
                if (!store.lat || !store.lng) return null;
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

        fetch("/aggregates")
            .then((res) => res.json())
            .then((data) => {
                setTotalOrders(data.all_store_total_order_count);
                setTotalAmount(data.all_store_total_order_amount);
            })
            .catch((err) => console.error("Failed to fetch aggregates:", err));
    }, [storeMap]);

    return (
        <>
            <div className="map-container" ref={mapContainerRef} />

            {showOverlay && (
                <div className="data-overlay">
                    <h3>Aggregates</h3>
                    <p>
                        <strong>Total Orders: </strong>{" "}
                        {totalOrders.toLocaleString()}
                    </p>
                    <p>
                        <strong>Total Amount: </strong> $
                        {totalAmount.toLocaleString()}
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
                    style={{
                        backgroundColor: "rgba(0,0,0,0.4)",
                    }}
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

export default App;
