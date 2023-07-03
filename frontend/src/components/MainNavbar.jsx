import Container from "react-bootstrap/Container";
import Nav from "react-bootstrap/Nav";
import Navbar from "react-bootstrap/Navbar";
import { Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

function MainNavbar() {
    const { logout, authUser, isLoggedIn } = useAuth();

    const handleLogout = async (event) => {
        logout();
    };

    return (
        <Navbar collapseOnSelect expand="lg" className="bg-body-tertiary">
            <Container>
                <Navbar.Brand as={Link} to="/">
                    My Top XYZ
                </Navbar.Brand>
                <Navbar.Toggle aria-controls="responsive-navbar-nav" />
                <Navbar.Collapse id="responsive-navbar-nav">
                    <Nav className="ms-auto">
                        {!isLoggedIn && (
                            <>
                                <Nav.Link as={Link} to="/register">
                                    Register
                                </Nav.Link>
                                <Nav.Link as={Link} to="/login">
                                    Login
                                </Nav.Link>
                            </>
                        )}
                        {isLoggedIn && (
                            <>
                                <Navbar.Text className="nav-link-style">
                                    {authUser.email}
                                </Navbar.Text>
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
