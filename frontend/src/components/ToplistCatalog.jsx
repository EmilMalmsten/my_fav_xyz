import { Link } from "react-router-dom";
import axios from "axios";
import { Row, Col } from "react-bootstrap";
import { useEffect, useState } from "react";
import PropTypes from "prop-types";
import defaultImage from "../assets/defaultItemImage.jpg";

function ToplistCatalog({ title, endpoint }) {
    const [toplists, setToplists] = useState([]);

    function getImageUrl(toplist) {
        if (toplist && toplist.items && Array.isArray(toplist.items)) {
            const sortedItems = toplist.items.sort((a, b) => a.rank - b.rank);
            const itemWithImagePath = sortedItems.find(
                (item) => item.image_path
            );

            if (itemWithImagePath) {
                return `http://localhost:8080/images/${itemWithImagePath.list_id}/${itemWithImagePath.image_path}`;
            }
        }
        return defaultImage;
    }

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await axios.get(
                    import.meta.env.VITE_API_URL + endpoint,
                    {
                        params: {
                            page_size: "8",
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
            <h5 className="my-3">{title}</h5>
            <Row>
                {toplists.map((toplist) => (
                    <Col xs={12} sm={6} md={4} lg={3}>
                        <Link
                            to={`/toplists/${toplist.toplist_id}`}
                            style={{ textDecoration: "none" }}
                        >
                            <div className="overlay-container">
                                <div
                                    style={{
                                        backgroundImage: `url(${getImageUrl(
                                            toplist
                                        )})`,
                                        backgroundSize: "cover",
                                        width: "100%",
                                        height: "100%",
                                        position: "absolute",
                                        zIndex: 1,
                                    }}
                                />
                                <div className="overlay" />
                                <div
                                    style={{
                                        position: "relative",
                                        zIndex: 3,
                                        color: "white",
                                        width: "100%",
                                        height: "100%",
                                        display: "flex",
                                        padding: "0 10px",
                                        fontSize: "20px",
                                    }}
                                >
                                    {toplist.title}
                                </div>
                            </div>
                        </Link>
                    </Col>
                ))}
            </Row>
        </>
    );
}

ToplistCatalog.propTypes = {
    title: PropTypes.string.isRequired,
    endpoint: PropTypes.string.isRequired,
};

export default ToplistCatalog;
