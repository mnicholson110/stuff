
import React, { useEffect, useState } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'leaflet/dist/leaflet.css';

// React Leaflet imports
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import L from 'leaflet';

// Hack to fix default marker icons not working with React Leaflet + Webpack
import iconUrl from 'leaflet/dist/images/marker-icon.png';
import iconRetinaUrl from 'leaflet/dist/images/marker-icon-2x.png';
import shadowUrl from 'leaflet/dist/images/marker-shadow.png';
L.Marker.prototype.options.icon = L.icon({
  iconRetinaUrl,
  iconUrl,
  shadowUrl,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41],
});

function App() {
  // We'll store store data in an object: { [storeId]: parsedJson }
  // Each parsedJson might have lat, long, orderCount, etc.
  const [storeMap, setStoreMap] = useState({});

  useEffect(() => {
    // 1) Fetch the snapshot on load
    fetch('http://localhost:8081/snapshot')
      .then((res) => res.json())
      .then((initialData) => {
        // initialData is { "2": "{\"storeId\":2,\"orderCount\":7,...}", "3": "{\"storeId\":3,...}", ...}
        // We need to parse each raw JSON string into an object
        const parsed = {};
        Object.entries(initialData).forEach(([storeId, rawJsonString]) => {
          try {
            const obj = JSON.parse(rawJsonString);
            parsed[storeId] = obj;
          } catch (e) {
            console.error('Error parsing store JSON:', e);
          }
        });
        setStoreMap(parsed);
      })
      .catch((err) => console.error('Snapshot fetch error:', err));

    // 2) Open SSE to /events
    const es = new EventSource('http://localhost:8081/events');
    es.onmessage = (event) => {
      // event.data is a raw JSON string for a single store update
      try {
        const updated = JSON.parse(event.data); 
        // e.g. { storeId: 2, orderCount: 10, storeLat: 40.7, storeLong: -73.9, ... }

        setStoreMap((prev) => {
          return {
            ...prev,
            [updated.storeId]: updated
          };
        });
      } catch (e) {
        console.error('SSE parse error:', e);
      }
    };

    // Cleanup on unmount
    return () => {
      es.close();
    };
  }, []);

  // We'll define a default map center if we have no data yet
  const defaultCenter = [40.7, -73.9]; // e.g., NYC
  const zoomLevel = 11;

  // Build an array of store data from the storeMap
  const storeData = Object.values(storeMap);

  return (
    <div className="container-fluid">
      <h1 className="my-3 text-center">Real-Time Store Dashboard</h1>
      <div className="row">
        <div className="col-12">
          <MapContainer
            center={defaultCenter}
            zoom={zoomLevel}
            style={{ height: '600px', width: '100%' }}
          >
            <TileLayer
              attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a>'
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            {storeData.map((store) => {
              // store might have storeLat, storeLong, storeId, orderCount, ...
              const lat = store.storeLat || 40.7;
              const lng = store.storeLong || -73.9;
              return (
                <Marker key={store.storeId} position={[lat, lng]}>
                  <Popup>
                    <div>
                      <h5>Store: {store.storeId}</h5>
                      <p>Order Count: {store.orderCount}</p>
                      {/* If you have more fields, show them here */}
                    </div>
                  </Popup>
                </Marker>
              );
            })}
          </MapContainer>
        </div>
      </div>
    </div>
  );
}

export default App;

