import { useLocation } from "react-router-dom";
import { Form, Button, Container } from "react-bootstrap";

function EditToplist() {
    const location = useLocation();
    const toplist = location.state || {};

    const handleSave = () => {
        // Handler function for the save button
        // You can implement the logic for saving the changes here
    };

    const handleCancel = () => {
        // Handler function for the cancel button
        // You can implement the logic for canceling the edit here
    };

    return (
        <Container style={{ maxWidth: "50%", margin: "3rem auto" }}>
            <h1>Edit Toplist</h1>
            <Form>
                <Form.Group controlId="formTitle">
                    <Form.Label>Title</Form.Label>
                    <Form.Control type="text" defaultValue={toplist.title} />
                </Form.Group>
                <Form.Group controlId="formDescription">
                    <Form.Label>Description</Form.Label>
                    <Form.Control
                        as="textarea"
                        rows={3}
                        defaultValue={toplist.description}
                    />
                </Form.Group>
                <Button variant="primary" onClick={handleSave}>
                    Save
                </Button>{" "}
                <Button variant="secondary" onClick={handleCancel}>
                    Cancel
                </Button>
            </Form>
        </Container>
    );
}

export default EditToplist;
