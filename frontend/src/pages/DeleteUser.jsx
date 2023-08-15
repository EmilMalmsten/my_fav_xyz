import { useState } from "react";
import { Form, Button, Container, Alert } from "react-bootstrap";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

function DeleteUser() {
    const [password, setPassword] = useState("");
    const [showFailureAlert, setShowFailureAlert] = useState(false);
    const [failureAlertMessage, setFailureAlertMessage] = useState("");
    const navigate = useNavigate();
    const { authUser, logout } = useAuth();

    const handleCancel = () => {
        navigate(`/users/${authUser.userID}`);
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        const form = event.target;
        if (form.checkValidity()) {
            try {
                await axios.delete(`${import.meta.env.VITE_API_URL}/users`, {
                    data: {
                        email: authUser.email,
                        password: password,
                    },
                });
                await logout();
                navigate("/", {
                    state: { successAlert: "Account deleted successfully" },
                });
            } catch (error) {
                console.error(error);
                if (error.response && error.response.data) {
                    setShowFailureAlert(true);
                    setFailureAlertMessage(error.response.data.error);
                } else {
                    setShowFailureAlert(true);
                    setFailureAlertMessage(
                        "Account deletion failed. Please try again."
                    );
                }
            }
        }
        form.classList.add("was-validated");
    };

    const handleInputChange = (event, setState) => {
        setState(event.target.value);
    };

    return (
        <Container style={{ maxWidth: "50%", margin: "3rem auto" }}>
            <div style={{ margin: "0 auto" }}>
                {showFailureAlert && (
                    <Alert
                        variant="danger"
                        onClose={() => setShowFailureAlert(false)}
                        dismissible
                    >
                        {failureAlertMessage}
                    </Alert>
                )}
                <h2>Enter password to delete account</h2>
                <Form noValidate onSubmit={handleSubmit}>
                    <Form.Group
                        controlId="password"
                        style={{ marginBottom: "1rem" }}
                    >
                        <Form.Label>Enter password</Form.Label>
                        <Form.Control
                            required
                            type="password"
                            placeholder="Enter password"
                            minLength="8"
                            value={password}
                            onChange={(e) => handleInputChange(e, setPassword)}
                        />
                        <Form.Control.Feedback type="invalid">
                            Password must be at least 8 characters.
                        </Form.Control.Feedback>
                    </Form.Group>

                    <Button
                        variant="outline-danger"
                        type="submit"
                        onClick={handleSubmit}
                    >
                        Delete Account
                    </Button>
                    <Button
                        variant="outline-secondary"
                        onClick={handleCancel}
                        className="mx-2"
                    >
                        Cancel
                    </Button>
                </Form>
            </div>
        </Container>
    );
}

export default DeleteUser;
