import React, { useEffect, useState, useRef } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'leaflet/dist/leaflet.css';
import { MapContainer, TileLayer, Marker, Popup, CircleMarker } from 'react-leaflet';
import L from 'leaflet';

function App() {
  const [storeMap, setStoreMap] = useState({});
  const [totalOrders, setTotalOrders] = useState(0);
  const [totalAmount, setTotalAmount] = useState(0);
  const [topStores, setTopStores] = useState([]);
  const [updatedStoreId, setUpdatedStoreId] = useState(null);
  const eventSourceRef = useRef(null);

  const formatNumber = (num) => num.toLocaleString();

  useEffect(() => {
    fetch('/snapshot')
      .then((res) => res.json())
      .then((data) => {
        const parsedData = {};
        let orders = 0;
        let amount = 0;

        Object.entries(data).forEach(([key, value]) => {
          const store = value;
          parsedData[key] = store;
          orders += store.order_count;
          amount += store.total_order_amount;
        });

        setStoreMap(parsedData);
        setTotalOrders(orders);
        setTotalAmount(amount);
        updateTopStores(parsedData);
      })
      .catch((err) => console.error('Failed to fetch snapshot:', err));

    const es = new EventSource('/events');
    eventSourceRef.current = es;

    es.onmessage = (event) => {
      try {
        const updatedStore = JSON.parse(event.data);
        setUpdatedStoreId(updatedStore.store_id);

        setStoreMap((prev) => {
          const updatedMap = { ...prev, [updatedStore.store_id]: updatedStore };

          const orders = Object.values(updatedMap).reduce((sum, store) => sum + store.order_count, 0);
          const amount = Object.values(updatedMap).reduce((sum, store) => sum + store.total_order_amount, 0);
          setTotalOrders(orders);
          setTotalAmount(amount);
          updateTopStores(updatedMap);

          setTimeout(() => setUpdatedStoreId(null), 1000); // Reset updated store ID after 1 second

          return updatedMap;
        });
      } catch (err) {
        console.error('SSE parse error:', err, event.data);
      }
    };

    return () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
      }
    };
  }, []);

  const updateTopStores = (storeData) => {
    const sortedStores = Object.values(storeData)
      .sort((a, b) => b.total_order_amount - a.total_order_amount)
      .slice(0, 5);
    setTopStores(sortedStores);
  };

  const defaultCenter = [38.897957, -77.03656];
  const defaultZoom = 5;

  return (
    <div className="container-fluid">
      <header className="bg-dark text-white py-3">
        <h1 className="text-center">Real-Time Store Dashboard</h1>
      </header>
      <div className="row my-3">
        <div className="col-md-6">
          <div className="card text-center shadow-sm bg-dark text-white">
            <div className="card-body">
              <h4 className="card-title">Total Orders</h4>
              <p className="card-text display-5">{formatNumber(totalOrders)}</p>
            </div>
          </div>
        </div>
        <div className="col-md-6">
          <div className="card text-center shadow-sm bg-dark text-white">
            <div className="card-body">
              <h4 className="card-title">Total Amount</h4>
              <p className="card-text display-5">${formatNumber(totalAmount)}</p>
            </div>
          </div>
        </div>
      </div>
      <div className="row my-3">
        <div className="col-md-6">
          <MapContainer center={defaultCenter} zoom={defaultZoom} style={{ height: '400px', width: '100%' }}>
            <TileLayer
              attribution="&copy; OpenStreetMap contributors"
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            {Object.values(storeMap).map((store) => {
              const lat = store.lat;
              const lng = store.long;
              if (!lat || !lng) return null;

              return (
                <CircleMarker
                  key={store.store_id}
                  center={[lat, lng]}
                  radius={5}
                  fillColor={store.store_id === updatedStoreId ? 'green' : 'blue'}
                  color={store.store_id === updatedStoreId ? 'green' : 'blue'}
                  fillOpacity={1}
                >
                  <Popup>
                    <div>
                      <strong>Store ID:</strong> {store.store_id} <br />
                      <strong>Order Count:</strong> {formatNumber(store.order_count)} <br />
                      <strong>Total Order Amount:</strong> ${formatNumber(store.total_order_amount)} <br />
                    </div>
                  </Popup>
                </CircleMarker>
              );
            })}
          </MapContainer>
        </div>
        <div className="col-md-6">
          <div className="card shadow-sm bg-dark text-white">
            <div className="card-body">
              <h4 className="card-title text-center">Top 5 Stores by Total Amount</h4>
              <table className="table table-dark table-striped">
                <thead>
                  <tr>
                    <th>Store ID</th>
                    <th>Order Count</th>
                    <th>Total Amount</th>
                  </tr>
                </thead>
                <tbody>
                  {topStores.map((store) => (
                    <tr key={store.store_id}>
                      <td>{store.store_id}</td>
                      <td>{formatNumber(store.order_count)}</td>
                      <td>${formatNumber(store.total_order_amount)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
