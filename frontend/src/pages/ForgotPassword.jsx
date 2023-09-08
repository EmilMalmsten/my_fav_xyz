import { useState } from "react";
import { Form, Button, Container, Alert } from "react-bootstrap";
import axios from "axios";
import { useNavigate } from "react-router-dom";

function ForgotPassword() {
    const [email, setEmail] = useState("");
    const [showFailureAlert, setShowFailureAlert] = useState(false);
    const [failureAlertMessage, setFailureAlertMessage] = useState("");
    const navigate = useNavigate();

    const handleSubmit = async (event) => {
        event.preventDefault();
        const form = event.target;
        if (form.checkValidity()) {
            try {
                const resp = await axios.post(
                    `${import.meta.env.VITE_API_URL}/forgotpassword`,
                    {
                        email: email,
                    }
                );
                console.log(resp);
                navigate("/", {
                    state: { successAlert: resp.data },
                });
            } catch (error) {
                if ((error.response.status = 404)) {
                    setShowFailureAlert(true);
                    setFailureAlertMessage(
                        "No user found with that email address"
                    );
                }
                console.log(error);
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
                <h4>Forgot password</h4>
                <Form noValidate onSubmit={handleSubmit}>
                    <Form.Group
                        controlId="email"
                        style={{ marginBottom: "1rem" }}
                    >
                        <Form.Label>
                            Enter email connected to your account
                        </Form.Label>
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

                    <Button variant="outline-primary" type="submit">
                        Reset Password
                    </Button>
                </Form>
            </div>
        </Container>
    );
}

export default ForgotPassword;
