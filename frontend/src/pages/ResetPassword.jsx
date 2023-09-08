import { useState } from "react";
import { Form, Button, Container, Alert } from "react-bootstrap";
import axios from "axios";
import { useNavigate, useParams } from "react-router-dom";

function ResetPassword() {
    const [newPassword, setNewPassword] = useState("");
    const [repeatNewPassword, setRepeatNewPassword] = useState("");
    const [showFailureAlert, setShowFailureAlert] = useState(false);
    const [failureAlertMessage, setFailureAlertMessage] = useState("");
    let { resetToken } = useParams();
    const navigate = useNavigate();

    const handleSubmit = async (event) => {
        event.preventDefault();
        const form = event.target;
        if (form.checkValidity()) {
            try {
                await axios.patch(
                    `${
                        import.meta.env.VITE_API_URL
                    }/resetpassword/${resetToken}`,
                    {
                        password: newPassword,
                    }
                );
                navigate("/", {
                    state: { successAlert: "Password updated successfully" },
                });
            } catch (error) {
                setShowFailureAlert(true);
                setFailureAlertMessage(
                    "Password update failed. Please contact support"
                );
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
                        controlId="newPassword"
                        style={{ marginBottom: "1rem" }}
                    >
                        <Form.Label>Enter New Password</Form.Label>
                        <Form.Control
                            required
                            type="password"
                            placeholder="Enter new password"
                            minLength="8"
                            value={newPassword}
                            onChange={(e) =>
                                handleInputChange(e, setNewPassword)
                            }
                        />
                        <Form.Control.Feedback type="invalid">
                            Password must be at least 8 characters.
                        </Form.Control.Feedback>
                    </Form.Group>

                    <Form.Group
                        controlId="repeatNewPassword"
                        style={{ marginBottom: "1rem" }}
                    >
                        <Form.Label>Repeat New Password</Form.Label>
                        <Form.Control
                            required
                            type="password"
                            placeholder="Repeat new password"
                            minLength="8"
                            value={repeatNewPassword}
                            onChange={(e) =>
                                handleInputChange(e, setRepeatNewPassword)
                            }
                            pattern={`^${newPassword}$`}
                        />
                        <Form.Control.Feedback type="invalid">
                            Passwords do not match.
                        </Form.Control.Feedback>
                    </Form.Group>

                    <Button variant="outline-primary" type="submit">
                        Save
                    </Button>
                </Form>
            </div>
        </Container>
    );
}

export default ResetPassword;
