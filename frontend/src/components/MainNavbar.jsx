import Container from "react-bootstrap/Container";
import Nav from "react-bootstrap/Nav";
import Navbar from "react-bootstrap/Navbar";
import Button from "react-bootstrap/Button";
import logo from "../assets/logo.svg";
import { Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";

function MainNavbar() {
    const { logout, authUser, isLoggedIn } = useAuth();
    const navigate = useNavigate();

    const handleLogout = async () => {
        logout();
    };

    const handleViewUser = () => {
        navigate(`/users/${authUser.userID}`);
    };

    return (
        <Navbar
            collapseOnSelect
            expand="lg"
            style={{
                borderBottom: "1px solid #ddd",
                boxShadow: "0 2px 2px -2px rgba(0,0,0,.1)",
            }}
        >
            <Container>
                <Navbar.Brand as={Link} to="/">
                    <img src={logo} style={{ width: "300px" }}></img>
                </Navbar.Brand>
                <Navbar.Toggle aria-controls="responsive-navbar-nav" />
                <Navbar.Collapse id="responsive-navbar-nav">
                    <Nav className="ms-auto">
                        {!isLoggedIn && (
                            <>
                                <Button
                                    as={Link}
                                    to="/login"
                                    variant="outline-secondary"
                                    className="mb-2 mb-lg-0 me-lg-3"
                                >
                                    Log in
                                </Button>
                                <Button
                                    as={Link}
                                    to="/register"
                                    className="brand-button"
                                >
                                    Sign up
                                </Button>
                            </>
                        )}
                        {isLoggedIn && (
                            <>
                                <Nav.Link onClick={handleViewUser}>
                                    {authUser.email}
                                </Nav.Link>
                                <Nav.Link onClick={handleLogout}>
                                    Logout
                                </Nav.Link>
                            </>
                        )}
                    </Nav>
                </Navbar.Collapse>
            </Container>
        </Navbar>
    );
}

export default MainNavbar;
