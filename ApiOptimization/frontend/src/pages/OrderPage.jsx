import { useParams, useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";

export default function OrderPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [order, setOrder] = useState(null);
  const [error, setError] = useState("");

  useEffect(() => {
    fetch(`http://localhost:8080/order/${id}`)
      .then((res) => {
        if (!res.ok) throw new Error("Заказ не найден");
        return res.json();
      })
      .then(setOrder)
      .catch((err) => setError(err.message));
  }, [id]);

  if (error) return <div style={{ padding: "20px" }}>Ошибка: {error}</div>;
  if (!order) return <div style={{ padding: "20px" }}>Загрузка...</div>;

  return (
    <div style={{ padding: "20px" }}>
      <button
        onClick={() => navigate("/")}
        style={{
          padding: "8px 12px",
          marginBottom: "20px",
          background: "#ddd",
          border: "1px solid #ccc",
          cursor: "pointer",
        }}
      >
        ⬅ Назад к поиску
      </button>

      <h1>Заказ {order.order_uid}</h1>
      <p><strong>Трек-номер:</strong> {order.track_number}</p>
      <p><strong>Покупатель:</strong> {order.customer_id}</p>
      <p><strong>Служба доставки:</strong> {order.delivery_service}</p>
      <p><strong>Дата создания:</strong> {new Date(order.date_created).toLocaleString()}</p>
      <p><strong>Сумма:</strong> {order.total_price} ₽</p>
      <h3>Товары:</h3>
      <ul>
        {order.items.map((item, idx) => (
          <li key={idx}>{item}</li>
        ))}
      </ul>
    </div>
  );
}
