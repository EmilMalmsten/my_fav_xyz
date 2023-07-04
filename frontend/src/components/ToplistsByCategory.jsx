import { Link } from "react-router-dom";
import axios from "axios";
import { useEffect, useState } from "react";
import PropTypes from "prop-types";

function ToplistsByCategory({ title, endpoint }) {
    const [toplists, setToplists] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await axios.get(
                    import.meta.env.VITE_API_URL + endpoint,
                    {
                        params: {
                            page_size: "20",
                        },
                    }
                );
                console.log(response.data);
                setToplists(response.data);
            } catch (error) {
                console.error(error);
            }
        };

        fetchData();
    }, []);

    return (
        <>
            <h5>{title}</h5>
            <ul>
                {toplists.map((item) => (
                    <li key={item.toplist_id}>
                        <Link to={`/toplists/${item.toplist_id}`}>
                            {item.title}
                        </Link>
                    </li>
                ))}
            </ul>
        </>
    );
}

ToplistsByCategory.propTypes = {
    title: PropTypes.string.isRequired,
    endpoint: PropTypes.string.isRequired,
};

export default ToplistsByCategory;
