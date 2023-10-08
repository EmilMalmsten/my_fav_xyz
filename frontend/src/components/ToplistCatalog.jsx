import { useNavigate } from "react-router-dom";
import axios from "axios";
import { Row, Col, Container } from "react-bootstrap";
import { useEffect, useState } from "react";
import PropTypes from "prop-types";
import defaultImage from "../assets/defaultItemImage.jpg";

function ToplistCatalog({ title, endpoint }) {
    const [toplists, setToplists] = useState([]);
    const navigate = useNavigate();

    function getImageUrl(toplist) {
        if (toplist && toplist.items && Array.isArray(toplist.items)) {
            const sortedItems = toplist.items.sort((a, b) => a.rank - b.rank);
            const itemWithImagePath = sortedItems.find(
                (item) => item.image_path
            );

            if (itemWithImagePath) {
                return `${import.meta.env.VITE_IMG_URL}/${
                    itemWithImagePath.list_id
                }/${itemWithImagePath.image_path}`;
            }
        }
        return defaultImage;
    }

    const handleClick = (toplist) => {
        navigate(`/toplists/${toplist.toplist_id}`);
    };

    useEffect(() => {
        const fetchData = async () => {
            try {
                const requestURL = import.meta.env.VITE_API_URL + endpoint;
                const response = await axios.get(requestURL, {
                    params: {
                        page_size: "12",
                    },
                });
                if (response.data) {
                    setToplists(response.data);
                }
            } catch (error) {
                console.error(error);
            }
        };

        fetchData();
    }, []);

    return (
        <>
            <h5 className="my-3">{title}</h5>

            <Row style={{ margin: "0 -5px" }}>
                {toplists.map((toplist) => (
                    <Col
                        key={toplist.created_at}
                        xs={12}
                        sm={6}
                        md={3}
                        lg={2}
                        style={{
                            display: "flex",
                            alignItems: "center",
                            justifyContent: "center",
                            padding: "5px 5px",
                        }}
                    >
                        <div
                            className="square-container"
                            onClick={() => handleClick(toplist)}
                            style={{
                                backgroundImage: `url(${getImageUrl(toplist)})`,
                                backgroundSize: "cover",
                            }}
                        >
                            <div className="overlay">
                                <p className="mx-2">{toplist.title}</p>
                            </div>
                        </div>
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
