import { useLocation, useNavigate } from "react-router-dom";
import { Container, Card, Form, Button, Row, Col } from "react-bootstrap";
import { useState } from "react";
import axios from "axios";

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

                // Add the item data
                items.forEach((item, index) => {
                    formData.append(`items[${index}][title]`, item.title);
                    formData.append(
                        `items[${index}][description]`,
                        item.description
                    );
                    formData.append(`items[${index}][rank]`, item.rank);
                    if (item.newImageFile) {
                        console.log("apending image");
                        formData.append(
                            `items[${index}][image]`,
                            item.newImageFile
                        );
                    }
                });

                await axios.put(
                    `${import.meta.env.VITE_API_URL}/toplists`,
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

    const renderItem = (item, index) => {
        let imageSource;
        if (item.newImageURL) {
            imageSource = item.newImageURL;
        } else if (item.image_path !== "") {
            imageSource = `http://localhost:8080/images/${item.list_id}/${item.image_path}`;
        }
        return (
            <Card
                key={index}
                style={{
                    marginBottom: "1rem",
                }}
            >
                <Row>
                    <Col xs={6} md={6}>
                        {imageSource && (
                            <img
                                src={imageSource}
                                alt={`Item ${index + 1}`}
                                style={{
                                    width: "80%",
                                    height: "60%",
                                    objectFit: "contain",
                                    marginBottom: "1rem",
                                }}
                            />
                        )}
                        <Button
                            onClick={() => {
                                const fileInput =
                                    document.createElement("input");
                                fileInput.type = "file";
                                fileInput.accept = "image/*";
                                fileInput.onchange = (e) => {
                                    if (e.target.files && e.target.files[0]) {
                                        handleImageUpload(
                                            index,
                                            e.target.files[0]
                                        );
                                    }
                                };
                                fileInput.click();
                            }}
                        >
                            {imageSource ? "Change Image" : "Upload Image"}
                        </Button>
                        {imageSource && (
                            <Button
                                variant="danger"
                                style={{ marginLeft: "1rem" }}
                                onClick={() => {
                                    handleItemChange(index, "image_path", "");
                                    handleItemChange(index, "newImageURL", "");
                                }}
                            >
                                Remove Image
                            </Button>
                        )}
                    </Col>
                    <Col xs={6} md={6}>
                        <Card.Body>
                            <Card.Title>Rank: {item.rank}</Card.Title>
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
                            <Button
                                variant="danger"
                                onClick={() => handleRemoveItem(index)}
                            >
                                Remove
                            </Button>
                        </Card.Body>
                    </Col>
                </Row>
            </Card>
        );
    };

    return (
        <Container style={{ maxWidth: "75%", margin: "3rem auto" }}>
            {items.map(renderItem)}
            <Button variant="primary" onClick={handleAddItem}>
                Add Item
            </Button>
            <Button variant="warning" onClick={handleCancel}>
                Cancel
            </Button>
            <Button variant="success" onClick={handleSave}>
                Save
            </Button>
        </Container>
    );
}

export default EditToplistItems;
