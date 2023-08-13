import { useState, useEffect } from "react";
import { Form, Button, Container, Alert } from "react-bootstrap";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

function UpdateUserEmail() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [showFailureAlert, setShowFailureAlert] = useState(false);
    const [failureAlertMessage, setFailureAlertMessage] = useState("");
    const navigate = useNavigate();
    const { authUser, updateUserInfo } = useAuth();

    const handleCancel = () => {
        navigate(`/users/${authUser.userID}`);
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        const form = event.target;
        if (form.checkValidity()) {
            try {
                const response = await axios.put(
                    `${import.meta.env.VITE_API_URL}/users/email`,
                    {
                        old_email: authUser.email,
                        new_email: email,
                        password: password,
                    }
                );
                updateUserInfo(
                    response.data.email,
                    response.data.id,
                    response.data.created_at
                );
                navigate("/", {
                    state: { successAlert: "Email updated successfully" },
                });
            } catch (error) {
                console.error(error);
                if (error.response && error.response.data) {
                    setShowFailureAlert(true);
                    setFailureAlertMessage(error.response.data.error);
                } else {
                    setShowFailureAlert(true);
                    setFailureAlertMessage(
                        "Email update failed. Please try again."
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
                <Form noValidate onSubmit={handleSubmit}>
                    <Form.Group
                        controlId="email"
                        style={{ marginBottom: "1rem" }}
                    >
                        <Form.Label>Enter New Email</Form.Label>
                        <Form.Control
                            required
                            type="email"
                            placeholder="Enter email"
                            value={email}
                            onChange={(e) => handleInputChange(e, setEmail)}
                        />
                        <Form.Control.Feedback type="invalid">
                            Please provide a valid email.
                        </Form.Control.Feedback>
                    </Form.Group>

                    <Form.Group
                        controlId="password"
                        style={{ marginBottom: "1rem" }}
                    >
                        <Form.Label>Enter Current Password</Form.Label>
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

                    <Button variant="outline-primary" type="submit">
                        Save
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

export default UpdateUserEmail;
