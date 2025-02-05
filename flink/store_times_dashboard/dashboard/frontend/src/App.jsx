import "maplibre-gl/dist/maplibre-gl.css";
import "./App.css";

import maplibregl from "maplibre-gl";
import React, { useEffect, useRef, useState } from "react";

const DARK_MAP_STYLE =
    "https://basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json";

function formatDuration(ms) {
    const seconds = ms / 1000;
    const formattedSeconds = seconds.toFixed(2);
    const unit = seconds === 1 ? "seconds" : "seconds";
    return `${formattedSeconds} ${unit}`;
}

export default function App() {
    const mapContainerRef = useRef(null);
    const mapRef = useRef(null);

    const [mapLoaded, setMapLoaded] = useState(true);
    const [storeMap, setStoreMap] = useState({});

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
                const { store_id, count, avg_processing, avg_delivered, avg_shipped } =
                    feature.properties;
                new maplibregl.Popup()
                    .setLngLat(e.lngLat)
                    .setHTML(
                        `
            <div style="color: #333;">
              <strong>Store ID:</strong> ${store_id}<br/>
              <strong>Order Count:</strong> ${count}<br/>
              <strong>Avg. Time to Processing: </strong> ${formatDuration(avg_processing)}<br/>
              <strong>Avg. Time to Shipped: </strong> ${formatDuration(avg_shipped)}<br/>
              <strong>Avg. Time to Delivered: </strong> ${formatDuration(avg_delivered)}<br/>
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
        }, 100);

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
                if (
                    !store.lat ||
                    !store.lng ||
                    !store.avg_shipped ||
                    !store.avg_delivered ||
                    !store.avg_processing
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
                        count: store.count,
                        avg_processing: store.avg_processing,
                        avg_delivered: store.avg_delivered,
                        avg_shipped: store.avg_shipped,
                    },
                };
            })
            .filter(Boolean);

        src.setData({
            type: "FeatureCollection",
            features,
        });
    }, [storeMap]);

    return (
        <>
            <div className="map-container" ref={mapContainerRef} />
        </>
    );
}
