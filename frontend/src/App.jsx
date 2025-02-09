import { useEffect, useState } from "react";
import { Table } from "antd";

const BACKEND_URL = "http://localhost:8080";

const formatDate = (isoString) => {
  if (!isoString) return "N/A";
  return new Intl.DateTimeFormat("ru-RU", {
    dateStyle: "long",
    timeStyle: "medium",
  }).format(new Date(isoString));
};

const columns = [
  { title: "IP адрес", dataIndex: "ip", key: "ip" },
  { 
    title: "Время пинга", 
    dataIndex: "pingtime", 
    key: "pingtime",
    render: (text) => formatDate(text) // Форматируем время пинга
  },
  { 
    title: "Последняя успешная попытка", 
    dataIndex: "successdate", 
    key: "successdate",
    render: (text) => formatDate(text) // Форматируем дату последней успешной попытки
  },
];

function App() {
  const [data, setData] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(`${BACKEND_URL}/containers`);
        const result = await response.json();
        setData(result);
      } catch (error) {
        console.error("Ошибка при загрузке данных:", error);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 5000); // Обновление каждые 5 сек.

    return () => clearInterval(interval);
  }, []);

  return (
    <div style={{ padding: "20px" }}>
      <h1>Мониторинг IP-адресов</h1>
      <Table dataSource={data} columns={columns} rowKey="ip" />
    </div>
  );
}

export default App;
