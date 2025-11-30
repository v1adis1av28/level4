import { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function SearchPage() {
  const [uuid, setUuid] = useState("");
  const navigate = useNavigate();

  const handleSearch = () => {
    if (uuid.trim()) {
      navigate(`/order/${uuid}`);
    }
  };

  return (
    <div style={{ padding: "20px" }}>
      <h1>Поиск заказа</h1>
      <input
        type="text"
        placeholder="Введите UUID заказа"
        value={uuid}
        onChange={(e) => setUuid(e.target.value)}
        style={{ padding: "8px", width: "300px" }}
      />
      <button onClick={handleSearch} style={{ padding: "8px 12px", marginLeft: "10px" }}>
        Найти
      </button>
    </div>
  );
}
