import { Container, Card, Form, Button } from "react-bootstrap";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const CreateToplist = () => {
    const navigate = useNavigate();

    const handleSubmit = async (event) => {
        event.preventDefault();
        const form = event.currentTarget;

        if (form.checkValidity()) {
            try {
                const accessToken = localStorage.getItem("accessToken");
                const response = await axios.post(
                    `${import.meta.env.VITE_API_URL}/toplists`,
                    {
                        title: form.elements.title.value,
                        description: form.elements.description.value,
                    },
                    {
                        headers: {
                            Authorization: `Bearer ${accessToken}`,
                        },
                    }
                );
                navigate(`/toplists/${response.data.id}`);
            } catch (error) {
                console.log(error);
            }
        }
        form.classList.add("was-validated");
    };

    return (
        <Container style={{ width: "50%" }}>
            <Card style={{ padding: "20px", margin: "20px" }}>
                <Card.Header>Create Toplist</Card.Header>
                <Card.Body>
                    <Form noValidate onSubmit={handleSubmit}>
                        <Form.Group
                            controlId="formTitle"
                            style={{ marginBottom: "15px" }}
                        >
                            <Form.Label>What is your toplist for?</Form.Label>
                            <Form.Control
                                type="text"
                                name="title"
                                placeholder="Enter title"
                                required
                                maxLength="100"
                            />
                            <Form.Control.Feedback type="invalid">
                                Please provide a title (up to 100 characters).
                            </Form.Control.Feedback>
                        </Form.Group>

                        <Form.Group
                            controlId="formDescription"
                            style={{ marginBottom: "15px" }}
                        >
                            <Form.Label>Description</Form.Label>
                            <Form.Control
                                as="textarea"
                                rows={3}
                                name="description"
                                placeholder="Enter description"
                                maxLength="1000"
                            />
                            <Form.Control.Feedback type="invalid">
                                Description can have up to 1000 characters.
                            </Form.Control.Feedback>
                        </Form.Group>

                        <Button
                            variant="primary"
                            type="submit"
                            className="brand-button"
                        >
                            Create
                        </Button>
                    </Form>
                </Card.Body>
            </Card>
        </Container>
    );
};

export default CreateToplist;
