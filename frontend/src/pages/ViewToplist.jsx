import React from "react";
import { useParams } from "react-router-dom";
import axios from "axios";
import { useEffect, useState, useRef } from "react";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";
import { Container, Col, Row, Button } from "react-bootstrap";
import ToplistItemImage from "../components/ToplistItemImage";

function Toplist() {
    const [toplist, setToplist] = useState({});
    const viewsUpdatedRef = useRef(false);
    const { authUser, isLoggedIn } = useAuth();
    const navigate = useNavigate();

    const { id } = useParams();

    const handleToplistEdit = () => {
        navigate(`/toplists/${id}/edit`, { state: toplist });
    };

    const handleToplistDelete = async () => {
        const confirmed = window.confirm(
            "Are you sure you want to delete the toplist?"
        );

        if (confirmed) {
            try {
                const accessToken = localStorage.getItem("accessToken");
                await axios.delete(
                    `${import.meta.env.VITE_API_URL}/toplists/${
                        toplist.toplist_id
                    }`,
                    {
                        headers: {
                            Authorization: `Bearer ${accessToken}`,
                        },
                    }
                );
                navigate("/", {
                    state: { successAlert: "Toplist deleted successfully" },
                });
            } catch (e) {
                console.error(e);
            }
        }
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
                console.log(response.data);
                setToplist(response.data);
            } catch (error) {
                console.error(error);
            }
        };

        if (!viewsUpdatedRef.current) {
            updateToplistViews();
        }
        fetchData();
    }, []);
    return (
        <>
            <Container style={{ width: "80%", margin: "2rem auto" }}>
                <div
                    className="my-4"
                    style={{ display: "flex", alignItems: "center" }}
                >
                    <h1>{toplist.title}</h1>

                    {isLoggedIn &&
                    Number(toplist.user_id) === Number(authUser.userID) ? (
                        <span className="mx-2" onClick={handleToplistEdit}>
                            ✏️
                        </span>
                    ) : null}
                </div>
                <p>{toplist.description}</p>
                {toplist.items && (
                    <>
                        {toplist.items.map((item) => (
                            <React.Fragment key={item.item_id}>
                                <Row className="my-5">
                                    <Col xs={1} s={1} md={1}>
                                        <h4>{item.rank}</h4>
                                    </Col>
                                    <Col
                                        xs={11}
                                        s={6}
                                        md={4}
                                        style={{ maxWidth: "220px" }}
                                    >
                                        <ToplistItemImage item={item} />
                                    </Col>
                                    <Col xs={12} s={5} md={7} className="mx-4">
                                        <h5>{item.title}</h5>
                                        <p>{item.description}</p>
                                    </Col>
                                </Row>
                                <hr />
                            </React.Fragment>
                        ))}
                    </>
                )}

                {isLoggedIn &&
                Number(toplist.user_id) === Number(authUser.userID) ? (
                    <div className="my-5">
                        {toplist.items === null ? (
                            <Button
                                variant="outline-primary"
                                onClick={handleToplistItems}
                            >
                                Add items
                            </Button>
                        ) : (
                            <Button
                                variant="outline-primary"
                                onClick={handleToplistItems}
                            >
                                Edit items
                            </Button>
                        )}

                        <Button
                            variant="outline-danger"
                            onClick={handleToplistDelete}
                            className="mx-2"
                        >
                            Delete Toplist
                        </Button>
                    </div>
                ) : null}
            </Container>
        </>
    );
}

export default Toplist;
