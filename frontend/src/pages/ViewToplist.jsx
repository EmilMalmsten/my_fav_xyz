import { useParams } from "react-router-dom";
import axios from "axios";
import { useEffect, useState, useRef } from "react";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";

function Toplist() {
    const [toplist, setToplist] = useState({});
    const viewsUpdatedRef = useRef(false);
    const { authUser, isLoggedIn } = useAuth();
    const navigate = useNavigate();

    const { id } = useParams();

    const handleToplistEdit = () => {
        navigate(`/toplists/${id}/edit`, { state: toplist });
    };

    const handleToplistItems = () => {
        navigate(`/toplists/${id}/items`, { state: toplist });
    };

    useEffect(() => {
        const updateToplistViews = async () => {
            await axios
                .post(`${import.meta.env.VITE_API_URL}/toplists/views/${id}`)
                .catch((error) => {
                    console.error(error);
                });
            viewsUpdatedRef.current = true;
        };

        const fetchData = async () => {
            try {
                const response = await axios.get(
                    `${import.meta.env.VITE_API_URL}/toplists/${id}`
                );
                setToplist(response.data);
                console.log(response.data);
            } catch (error) {
                console.error(error);
            }
        };

        if (!viewsUpdatedRef.current) {
            updateToplistViews();
        }
        fetchData();
        console.log(toplist);
    }, []);
    return (
        <>
            <div style={{ display: "flex", alignItems: "center" }}>
                <h1>{toplist.title}</h1>

                {isLoggedIn && toplist.userId === authUser.id ? (
                    <button onClick={handleToplistEdit}>Edit</button>
                ) : null}
            </div>
            <p>{toplist.description}</p>
            {toplist.items && (
                <ol>
                    {toplist.items.map((item) => (
                        <li key={item.item_id}>{item.title}</li>
                    ))}
                </ol>
            )}
            {toplist.items === null ? (
                <button onClick={handleToplistItems}>Add items</button>
            ) : (
                <button onClick={handleToplistItems}>Edit items</button>
            )}
        </>
    );
}

export default Toplist;
