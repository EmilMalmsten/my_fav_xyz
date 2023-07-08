import { Container, Card, Form, Button } from "react-bootstrap";

const CreateToplist = () => {
    const handleSubmit = (event) => {
        event.preventDefault();
        const form = event.currentTarget;

        if (form.checkValidity()) {
            console.log("Title:", form.elements.title.value);
            console.log("Description:", form.elements.description.value);

            // Reset form fields
            form.reset();
        } else {
            event.stopPropagation();
        }
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
                            <Form.Label>Title</Form.Label>
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

                        <Button variant="primary" type="submit">
                            Create
                        </Button>
                    </Form>
                </Card.Body>
            </Card>
        </Container>
    );
};

export default CreateToplist;
