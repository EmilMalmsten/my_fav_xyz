import { useParams } from "react-router-dom"
import axios from 'axios';
import React, { useEffect, useState } from 'react';

function Toplist() {
  const [toplist, setToplist] = useState({});

  const { id } = useParams();
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`${import.meta.env.VITE_API_URL}/toplists/${id}`);
        setToplist(response.data);
        console.log(response.data);
      } catch (error) {
        console.error(error);
      }
    };
  
    fetchData();
  }, []);
  return (
    <>
      <h1>{toplist.title}</h1>
      <p>{toplist.description}</p>
      {toplist.items && (
        <ol>
          {toplist.items.map((item) => (
            <li key={item.item_id}>{item.title}</li>
          ))}
        </ol>
      )}
    </>
  )
}

export default Toplist