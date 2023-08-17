import { useLocation, useNavigate } from "react-router-dom";
import { Form, Button, Container } from "react-bootstrap";
import axios from "axios";
import { useRef } from "react";

function EditToplist() {
    const location = useLocation();
    const navigate = useNavigate();
    const toplist = location.state || {};

    const formRef = useRef(null);

    const handleSave = async (event) => {
        event.preventDefault();
        const form = formRef.current;

        if (form.checkValidity()) {
            try {
                const accessToken = localStorage.getItem("accessToken");
                await axios.put(
                    `${import.meta.env.VITE_API_URL}/toplists`,
                    {
                        id: toplist.toplist_id,
                        title: form.elements.formTitle.value,
                        description: form.elements.formDescription.value,
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

    const handleCancel = () => {
        navigate(`/toplists/${toplist.toplist_id}`);
    };

    return (
        <Container style={{ maxWidth: "50%", margin: "3rem auto" }}>
            <h1>Edit Toplist</h1>
            <Form ref={formRef}>
                <Form.Group controlId="formTitle">
                    <Form.Label>Title</Form.Label>
                    <Form.Control
                        type="text"
                        defaultValue={toplist.title}
                        required
                    />
                </Form.Group>
                <Form.Group controlId="formDescription">
                    <Form.Label>Description</Form.Label>
                    <Form.Control
                        as="textarea"
                        rows={3}
                        defaultValue={toplist.description}
                        required
                    />
                </Form.Group>
                <Button variant="outline-primary" onClick={handleSave}>
                    Save
                </Button>
                <Button
                    className="m-2"
                    variant="outline-secondary"
                    onClick={handleCancel}
                >
                    Cancel
                </Button>
            </Form>
        </Container>
    );
}

export default EditToplist;
