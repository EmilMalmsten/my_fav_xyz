import axios from 'axios';
import Button from 'react-bootstrap/Button';
import React, { useEffect, useState } from 'react';

function ToplistCatalog({ title }) {
    const [items, setItems] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
          try {
            const response = await axios.get(`${import.meta.env.VITE_API_URL}/toplists`, {
              params: {
                page_id: '1',
                page_size: '5'
              }
            });
            setItems(response.data);
          } catch (error) {
            console.error(error);
          }
        };
      
        fetchData();
      }, []);

    return (
      <>
        <h5>{title}</h5>
        <ol>
            {items.map((item) => (
            <li key={item.id}>{item.title}</li>
            ))}
        </ol>
        <Button variant="outline-dark">View more</Button>
      </>
    )
}
  
export default ToplistCatalog