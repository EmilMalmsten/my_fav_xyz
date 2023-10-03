import { useLocation, useNavigate } from "react-router-dom";
import React from "react";
import { Container, Form, Button, Row, Col } from "react-bootstrap";
import { useState } from "react";
import axios from "axios";
import ToplistItemImage from "../components/ToplistItemImage";

function EditToplistItems() {
    const location = useLocation();
    const navigate = useNavigate();
    const toplist = location.state || {};

    const [items, setItems] = useState(
        toplist.items || [
            { title: "", description: "", rank: 1, image_path: "" },
        ]
    );

    const handleItemChange = (index, field, value) => {
        setItems((prevItems) =>
            prevItems.map((item, i) =>
                i === index
                    ? field === "image_path" && value === ""
                        ? {
                              ...item,
                              [field]: value,
                              newImageURL: null,
                              newImageFile: null,
                          }
                        : { ...item, [field]: value }
                    : item
            )
        );
    };

    const handleImageUpload = (index, file) => {
        const newItems = [...items];
        newItems[index] = {
            ...newItems[index],
            newImageURL: URL.createObjectURL(file),
            newImageFile: file,
        };
        setItems(newItems);
    };

    const handleAddItem = () => {
        const newRank = Math.max(...items.map((item) => item.rank)) + 1;
        const newItem = {
            title: "",
            description: "",
            rank: newRank,
            imageURL: "",
        };
        setItems((prevItems) => [...prevItems, newItem]);
    };

    const handleCancel = () => {
        navigate(`/toplists/${toplist.toplist_id}`);
    };

    const handleSave = async (event) => {
        event.preventDefault();
        const form = event.currentTarget;

        if (form.checkValidity()) {
            try {
                const accessToken = localStorage.getItem("accessToken");
                const formData = new FormData();

                // Add the toplist data
                formData.append("id", toplist.toplist_id);
                formData.append("title", toplist.title);
                formData.append("description", toplist.description);

                console.log(items);

                // Add the item data
                items.forEach((item, index) => {
                    formData.append(`items[${index}][title]`, item.title);
                    formData.append(
                        `items[${index}][description]`,
                        item.description
                    );
                    formData.append(`items[${index}][rank]`, item.rank);
                    if (item.image_path) {
                        formData.append(
                            `items[${index}][path]`,
                            item.image_path
                        );
                    }
                    if (item.newImageFile) {
                        formData.append(
                            `items[${index}][image]`,
                            item.newImageFile
                        );
                    }
                });

                await axios.put(
                    `${import.meta.env.VITE_API_URL}/toplists/items`,
                    formData,
                    {
                        headers: {
                            "Content-Type": "multipart/form-data",
                            Authorization: `Bearer ${accessToken}`,
                        },
                    }
                );
                navigate(`/toplists/${toplist.toplist_id}`);
            } catch (error) {
                console.log(error);
            }
        }
        form.classList.add("was-validated");
    };

    const handleRemoveItem = (index) => {
        setItems((prevItems) => {
            const newItems = prevItems.filter((item, i) => i !== index);
            return newItems.map((item) =>
                item.rank > prevItems[index].rank
                    ? { ...item, rank: item.rank - 1 }
                    : item
            );
        });
    };

    const handleMoveUp = (index) => {
        setItems((prevItems) => {
            if (index === 0) return prevItems;
            let newArray = [...prevItems];
            const temp = newArray[index];
            newArray[index] = newArray[index - 1];
            newArray[index - 1] = temp;

            // Update the rank of each item based on its new position
            newArray = newArray.map((item, i) => ({ ...item, rank: i + 1 }));

            return newArray;
        });
    };

    const handleMoveDown = (index) => {
        setItems((prevItems) => {
            if (index === prevItems.length - 1) return prevItems;
            let newArray = [...prevItems];
            const temp = newArray[index];
            newArray[index] = newArray[index + 1];
            newArray[index + 1] = temp;

            // Update the rank of each item based on its new position
            newArray = newArray.map((item, i) => ({ ...item, rank: i + 1 }));

            return newArray;
        });
    };
    const renderItem = (item, index) => {
        let imageSource;
        if (item.newImageURL) {
            imageSource = item.newImageURL;
        } else if (item.image_path) {
            imageSource = `http://localhost:8080/images/${item.list_id}/${item.image_path}`;
        }

        return (
            <React.Fragment key={item.item_id}>
                <Row>
                    <Col xs={1} s={1} md={1}>
                        <h4>{item.rank}</h4>
                    </Col>
                    <Col
                        xs={11}
                        s={6}
                        md={4}
                        className="py-4"
                        style={{ maxWidth: "220px" }}
                    >
                        <ToplistItemImage item={item} />
                        <div style={{ display: "flex" }}>
                            <span
                                onClick={() => {
                                    const fileInput =
                                        document.createElement("input");
                                    fileInput.type = "file";
                                    fileInput.accept = "image/*";
                                    fileInput.onchange = (e) => {
                                        if (
                                            e.target.files &&
                                            e.target.files[0]
                                        ) {
                                            handleImageUpload(
                                                index,
                                                e.target.files[0]
                                            );
                                        }
                                    };
                                    fileInput.click();
                                }}
                                className="emojiBtn"
                            >
                                ✏️
                            </span>
                            {imageSource && (
                                <span
                                    onClick={() => {
                                        handleItemChange(
                                            index,
                                            "image_path",
                                            ""
                                        );
                                        handleItemChange(
                                            index,
                                            "newImageURL",
                                            ""
                                        );
                                    }}
                                    className="emojiBtn"
                                >
                                    🗑️
                                </span>
                            )}
                        </div>
                    </Col>
                    <Col
                        xs={12}
                        s={5}
                        md={7}
                        style={{ display: "flex", alignItems: "flex-start" }}
                    >
                        <div style={{ flex: "1" }}>
                            <Form.Group controlId={`title-${index}`}>
                                <Form.Label>Title</Form.Label>
                                <Form.Control
                                    type="text"
                                    value={item.title}
                                    required
                                    maxLength="100"
                                    onChange={(e) =>
                                        handleItemChange(
                                            index,
                                            "title",
                                            e.target.value
                                        )
                                    }
                                />
                                <Form.Control.Feedback type="invalid">
                                    Please provide a title (up to 100
                                    characters).
                                </Form.Control.Feedback>
                            </Form.Group>
                            <Form.Group controlId={`description-${index}`}>
                                <Form.Label>Description</Form.Label>
                                <Form.Control
                                    as="textarea"
                                    rows={3}
                                    maxLength="1000"
                                    value={item.description}
                                    onChange={(e) =>
                                        handleItemChange(
                                            index,
                                            "description",
                                            e.target.value
                                        )
                                    }
                                />
                                <Form.Control.Feedback type="invalid">
                                    Description can have up to 1000 characters.
                                </Form.Control.Feedback>
                            </Form.Group>
                        </div>
                        <div
                            style={{
                                display: "flex",
                                flexDirection: "column",
                                alignItems: "flex-end",
                                justifyContent: "center",
                                height: "100%",
                                marginLeft: "10px",
                            }}
                        >
                            <span
                                className="emojiBtn"
                                onClick={() => handleRemoveItem(index)}
                            >
                                ❌
                            </span>
                            {index !== 0 && (
                                <span
                                    className="emojiBtn"
                                    onClick={() => handleMoveUp(index)}
                                >
                                    🔼
                                </span>
                            )}
                            {index !== items.length - 1 && (
                                <span
                                    className="emojiBtn"
                                    onClick={() => handleMoveDown(index)}
                                >
                                    🔽
                                </span>
                            )}
                        </div>
                    </Col>
                </Row>
                <hr />
            </React.Fragment>
        );
    };

    return (
        <Container style={{ width: "80%", margin: "3rem auto" }}>
            {items.map(renderItem)}
            <Button
                className="m-2"
                variant="outline-primary"
                onClick={handleAddItem}
            >
                Add Item
            </Button>
            <Button
                className="m-2"
                variant="outline-secondary"
                onClick={handleCancel}
            >
                Cancel Edit
            </Button>
            <Button
                className="m-2"
                variant="outline-success"
                onClick={handleSave}
            >
                Save Changes
            </Button>
        </Container>
    );
}

export default EditToplistItems;
