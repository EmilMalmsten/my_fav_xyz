import { useState } from "react";
import Container from "react-bootstrap/Container";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const CreateToplist = () => {
    const navigate = useNavigate();
    const [items, setItems] = useState([]);
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [rank, setRank] = useState(1);
    const [editItemId, setEditItemId] = useState(null);
    const [editTitle, setEditTitle] = useState("");
    const [editDescription, setEditDescription] = useState("");
    const [toplistTitle, setToplistTitle] = useState("");
    const [toplistDescription, setToplistDescription] = useState("");

    const addItem = () => {
        if (title.trim() !== "" && description.trim() !== "") {
            setItems((prevItems) => [
                ...prevItems,
                { rank: rank, title: title, description: description },
            ]);
            setTitle("");
            setDescription("");
            setRank((prevCounter) => prevCounter + 1);
        }
    };

    const deleteItem = (id) => {
        setItems((prevItems) => prevItems.filter((item) => item.id !== id));
    };

    const editItem = (id, title, description) => {
        setEditItemId(id);
        setEditTitle(title);
        setEditDescription(description);
    };

    const saveItem = () => {
        const updatedItems = items.map((item) => {
            if (item.id === editItemId) {
                return {
                    ...item,
                    title: editTitle,
                    description: editDescription,
                };
            }
            return item;
        });
        setItems(updatedItems);
        setEditItemId(null);
        setEditTitle("");
        setEditDescription("");
    };

    const cancelEdit = () => {
        setEditItemId(null);
        setEditTitle("");
        setEditDescription("");
    };

    const handleToplistCreate = async () => {
        const newToplist = {
            title: toplistTitle,
            description: toplistDescription,
            items: items,
        };

        try {
            const accessToken = localStorage.getItem("accessToken");
            const resp = await axios.post(
                `${import.meta.env.VITE_API_URL}/toplists`,
                newToplist,
                {
                    headers: {
                        Authorization: `Bearer ${accessToken}`,
                    },
                }
            );
            const toplistID = resp.data.id;
            navigate(`/toplists/${toplistID}`);
        } catch (error) {
            console.error(error);
        }
    };

    return (
        <Container style={{ width: "50%" }}>
            <Form>
                <Form.Group controlId="formToplistTitle">
                    <Form.Label>Toplist title:</Form.Label>
                    <Form.Control
                        type="text"
                        placeholder="Enter the title for your toplist"
                        value={toplistTitle}
                        onChange={(e) => setToplistTitle(e.target.value)}
                    />
                </Form.Group>
                <Form.Group controlId="formToplistDescription">
                    <Form.Label>Toplist description:</Form.Label>
                    <Form.Control
                        as="textarea"
                        rows={3}
                        placeholder="What is your toplist about?"
                        value={toplistDescription}
                        onChange={(e) => setToplistDescription(e.target.value)}
                    />
                </Form.Group>
            </Form>
            <ol>
                {items.map((item) => (
                    <li key={item.rank}>
                        {item.rank === editItemId ? (
                            <div>
                                <Form.Group controlId={`editTitle${item.rank}`}>
                                    <Form.Control
                                        type="text"
                                        placeholder="Enter title"
                                        value={editTitle}
                                        onChange={(e) =>
                                            setEditTitle(e.target.value)
                                        }
                                    />
                                </Form.Group>
                                <Form.Group
                                    controlId={`editDescription${item.rank}`}
                                >
                                    <Form.Control
                                        type="text"
                                        placeholder="Enter description"
                                        value={editDescription}
                                        onChange={(e) =>
                                            setEditDescription(e.target.value)
                                        }
                                    />
                                </Form.Group>
                                <Button
                                    variant="success"
                                    onClick={saveItem}
                                    className="mr-2"
                                >
                                    Save
                                </Button>
                                <Button
                                    variant="secondary"
                                    onClick={cancelEdit}
                                >
                                    Cancel
                                </Button>
                            </div>
                        ) : (
                            <div>
                                <h5>{item.title}</h5>
                                <p>{item.description}</p>
                                <Button
                                    variant="warning"
                                    size="sm"
                                    className="mr-2"
                                    onClick={() =>
                                        editItem(
                                            item.rank,
                                            item.title,
                                            item.description
                                        )
                                    }
                                >
                                    Edit
                                </Button>
                                <Button
                                    variant="danger"
                                    size="sm"
                                    onClick={() => deleteItem(item.rank)}
                                >
                                    Delete
                                </Button>
                            </div>
                        )}
                    </li>
                ))}
            </ol>
            <br />
            <Form>
                <Form.Label>Add items to your toplist</Form.Label>
                <Form.Group controlId="formNewTitle" className="d-flex">
                    <Form.Label className="mr-2">Enter title:</Form.Label>
                    <Form.Control
                        type="text"
                        placeholder="Enter title"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                    />
                </Form.Group>
                <Form.Group controlId="formNewDescription" className="d-flex">
                    <Form.Label className="mr-2">Enter description:</Form.Label>
                    <Form.Control
                        type="text"
                        placeholder="Enter description"
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                    />
                </Form.Group>
                <Button variant="primary" onClick={addItem}>
                    Add
                </Button>
            </Form>

            <div
                style={{
                    position: "fixed",
                    left: 0,
                    bottom: 20,
                    width: "100%",
                    textAlign: "center",
                }}
            >
                <p>Finished? Submit your toplist here</p>
                <Button variant="primary" onClick={handleToplistCreate}>
                    Create
                </Button>
            </div>
        </Container>
    );
};

export default CreateToplist;
