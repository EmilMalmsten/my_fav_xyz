import ToplistCatalog from "../components/ToplistCatalog";
import axios from "axios";
import { Container } from "react-bootstrap";
import { useParams, useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";

function ViewToplistsByUser() {
    const { id } = useParams();
    const [user, setUser] = useState({});
    const navigate = useNavigate();

    async function getUserInfo() {
        try {
            const resp = await axios.get(
                `${import.meta.env.VITE_API_URL}/users/${id}`
            );
            setUser(resp.data);
        } catch (e) {
            console.error(e);
            navigate("/404");
        }
    }

    useEffect(() => {
        getUserInfo();
    }, []);

    return (
        <Container>
            {user.email ? (
                <ToplistCatalog
                    title={`Toplists by ${user.email}`}
                    endpoint={`/toplists/user/${id}`}
                />
            ) : (
                <p>User doesn't exist</p>
            )}
        </Container>
    );
}

export default ViewToplistsByUser;
