import { Container, Row, Col } from "react-bootstrap";
import { useAuth } from "../context/AuthContext";

function Footer() {
    const { authUser, isLoggedIn } = useAuth();

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
                        <p>info@mytopxyz.com</p>
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
                                    <a
                                        href="#link1"
                                        style={{
                                            color: "#fff",
                                            textDecoration: "none",
                                        }}
                                    >
                                        My account
                                    </a>
                                </li>
                                <li
                                    style={{
                                        marginBottom: "10px",
                                        textAlign: "right",
                                    }}
                                >
                                    <a
                                        href="#link1"
                                        style={{
                                            color: "#fff",
                                            textDecoration: "none",
                                        }}
                                    >
                                        My toplists
                                    </a>
                                </li>
                                <li
                                    style={{
                                        marginBottom: "10px",
                                        textAlign: "right",
                                    }}
                                >
                                    <a
                                        href="#link1"
                                        style={{
                                            color: "#fff",
                                            textDecoration: "none",
                                        }}
                                    >
                                        Create new toplist
                                    </a>
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
                                    <a
                                        href="#link1"
                                        style={{
                                            color: "#fff",
                                            textDecoration: "none",
                                        }}
                                    >
                                        Login
                                    </a>
                                </li>
                                <li
                                    style={{
                                        marginBottom: "10px",
                                        textAlign: "right",
                                    }}
                                >
                                    <a
                                        href="#link1"
                                        style={{
                                            color: "#fff",
                                            textDecoration: "none",
                                        }}
                                    >
                                        Create Account
                                    </a>
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
                                <a
                                    href="#link1"
                                    style={{
                                        color: "#fff",
                                        textDecoration: "none",
                                    }}
                                >
                                    Cookie Policy
                                </a>
                            </li>
                            <li
                                style={{
                                    marginBottom: "10px",
                                    textAlign: "right",
                                }}
                            >
                                <a
                                    href="#link1"
                                    style={{
                                        color: "#fff",
                                        textDecoration: "none",
                                    }}
                                >
                                    Terms of Service
                                </a>
                            </li>
                            <li
                                style={{
                                    marginBottom: "10px",
                                    textAlign: "right",
                                }}
                            >
                                <a
                                    href="#link1"
                                    style={{
                                        color: "#fff",
                                        textDecoration: "none",
                                    }}
                                >
                                    Privacy Policy
                                </a>
                            </li>
                        </ul>
                    </Col>
                </Row>
            </Container>
        </div>
    );
}

export default Footer;
