import { useLocation } from "react-router-dom";
import { Container, Card, Form, Button } from "react-bootstrap";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

function EditToplistItems() {
    const location = useLocation();
    const navigate = useNavigate();
    const toplist = location.state || {};
    const [items, setItems] = useState(
        toplist.items || [{ title: "", description: "", rank: 1 }]
    );

    const handleItemChange = (index, field, value) => {
        setItems((prevItems) =>
            prevItems.map((item, i) =>
                i === index ? { ...item, [field]: value } : item
            )
        );
    };

    const handleAddItem = () => {
        const newRank = Math.max(...items.map((item) => item.rank)) + 1;
        const newItem = { title: "", description: "", rank: newRank };
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
                await axios.put(
                    `${import.meta.env.VITE_API_URL}/toplists`,
                    {
                        id: toplist.toplist_id,
                        title: toplist.title,
                        description: toplist.description,
                        items: items,
                    },
                    {
                        headers: {
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
        setItems((prevItems) => prevItems.filter((item, i) => i !== index));
    };

    const renderItem = (item, index) => {
        return (
            <Card key={index} style={{ marginBottom: "1rem" }}>
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
                                handleItemChange(index, "title", e.target.value)
                            }
                        />
                        <Form.Control.Feedback type="invalid">
                            Please provide a title (up to 100 characters).
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
            </Card>
        );
    };

    return (
        <Container style={{ maxWidth: "50%", margin: "3rem auto" }}>
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
