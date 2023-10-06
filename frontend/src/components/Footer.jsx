import { Container, Row, Col } from "react-bootstrap";
import { useAuth } from "../context/AuthContext";
import { Link } from "react-router-dom";

function Footer() {
    const { isLoggedIn, authUser } = useAuth();

    return (
        <div
            style={{
                backgroundColor: "#2E2E2E",
                padding: "50px 0",
                color: "#fff",
                marginTop: "auto",
            }}
        >
            <Container>
                <Row>
                    <Col>
                        <h5>My Top XYZ &copy; 2023</h5>
                        <p>topxyzinfo@gmail.com</p>
                    </Col>
                    <Col>
                        {isLoggedIn ? (
                            <ul style={{ listStyleType: "none", padding: "0" }}>
                                <li
                                    style={{
                                        marginBottom: "10px",
                                        textAlign: "right",
                                    }}
                                >
                                    <Link
                                        to={`/users/${authUser.userID}`}
                                        style={{
                                            color: "#fff",
                                            textDecoration: "none",
                                        }}
                                    >
                                        My account
                                    </Link>
                                </li>
                                <li
                                    style={{
                                        marginBottom: "10px",
                                        textAlign: "right",
                                    }}
                                >
                                    <Link
                                        to={`/toplists/user/${authUser.userID}`}
                                        style={{
                                            color: "#fff",
                                            textDecoration: "none",
                                        }}
                                    >
                                        My toplists
                                    </Link>
                                </li>
                                <li
                                    style={{
                                        marginBottom: "10px",
                                        textAlign: "right",
                                    }}
                                >
                                    <Link
                                        to="/toplists/create"
                                        style={{
                                            color: "#fff",
                                            textDecoration: "none",
                                        }}
                                    >
                                        Create new toplist
                                    </Link>
                                </li>
                            </ul>
                        ) : (
                            <ul style={{ listStyleType: "none", padding: "0" }}>
                                <li
                                    style={{
                                        marginBottom: "10px",
                                        textAlign: "right",
                                    }}
                                >
                                    <Link
                                        to="/login"
                                        style={{
                                            color: "#fff",
                                            textDecoration: "none",
                                        }}
                                    >
                                        Login
                                    </Link>
                                </li>
                                <li
                                    style={{
                                        marginBottom: "10px",
                                        textAlign: "right",
                                    }}
                                >
                                    <Link
                                        to="/register"
                                        style={{
                                            color: "#fff",
                                            textDecoration: "none",
                                        }}
                                    >
                                        Create Account
                                    </Link>
                                </li>
                            </ul>
                        )}
                    </Col>
                    <Col className="text-right">
                        <ul style={{ listStyleType: "none", padding: "0" }}>
                            <li
                                style={{
                                    marginBottom: "10px",
                                    textAlign: "right",
                                }}
                            >
                                <Link
                                    to="/cookies"
                                    style={{
                                        color: "#fff",
                                        textDecoration: "none",
                                    }}
                                >
                                    Cookie Policy
                                </Link>
                            </li>
                            <li
                                style={{
                                    marginBottom: "10px",
                                    textAlign: "right",
                                }}
                            >
                                <Link
                                    to="/tos"
                                    style={{
                                        color: "#fff",
                                        textDecoration: "none",
                                    }}
                                >
                                    Terms of Service
                                </Link>
                            </li>
                            <li
                                style={{
                                    marginBottom: "10px",
                                    textAlign: "right",
                                }}
                            >
                                <Link
                                    to="/privacy"
                                    style={{
                                        color: "#fff",
                                        textDecoration: "none",
                                    }}
                                >
                                    Privacy Policy
                                </Link>
                            </li>
                        </ul>
                    </Col>
                </Row>
            </Container>
        </div>
    );
}

export default Footer;
