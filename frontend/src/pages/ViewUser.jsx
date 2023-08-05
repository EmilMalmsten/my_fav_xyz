import { useAuth } from "../context/AuthContext";
import { useParams } from "react-router-dom";
import { Container } from "react-bootstrap";

function ViewUser() {
    const { authUser } = useAuth();
    const { id } = useParams();
    console.log(authUser, id);

    return (
        <Container style={{ maxWidth: "75%", margin: "3rem auto" }}>
            <p>Email: {authUser.email}</p>
            <p>User ID: {authUser.userID}</p>
        </Container>
    );
}

export default ViewUser;
