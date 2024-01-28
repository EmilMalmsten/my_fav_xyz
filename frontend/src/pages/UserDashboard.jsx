import { useAuth } from "../context/AuthContext";
import { useParams, useNavigate } from "react-router-dom";
import { Container, Row, Col, Button } from "react-bootstrap";

function UserDashboard() {
    const { authUser } = useAuth();
    const { id } = useParams();
    const navigate = useNavigate();

    const handleChangeEmail = () => {
        navigate(`/users/${authUser.userID}/email`);
    };

    const handleChangePassword = () => {
        navigate(`/users/${authUser.userID}/password`);
    };

    const handleDeleteAccount = () => {
        navigate(`/users/${authUser.userID}/delete`);
    };

    return (
        <Container className="text-center my-3">
            <h2>Manage account</h2>
            <Row className="justify-content-center">
                <Col xs={12} md={6}>
                    <Button
                        className="m-2 w-50 brand-button-outline"
                        variant="outline-primary"
                        onClick={handleChangeEmail}
                    >
                        Change Email Address
                    </Button>
                </Col>
            </Row>
            <Row className="justify-content-center">
                <Col xs={12} md={6}>
                    <Button
                        className="m-2 w-50 brand-button-outline"
                        variant="outline-primary"
                        onClick={handleChangePassword}
                    >
                        Change Password
                    </Button>
                </Col>
            </Row>
            <Row className="justify-content-center">
                <Col xs={12} md={6}>
                    <Button
                        className="m-2 w-50"
                        variant="outline-danger"
                        onClick={handleDeleteAccount}
                    >
                        Delete Account
                    </Button>
                </Col>
            </Row>
        </Container>
    );
}

export default UserDashboard;
