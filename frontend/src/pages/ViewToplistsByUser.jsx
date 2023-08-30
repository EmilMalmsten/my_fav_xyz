import ToplistCatalog from "../components/ToplistCatalog";
import { Container } from "react-bootstrap";
import { useParams } from "react-router-dom";

function ViewToplistsByUser() {
    const { id } = useParams();
    console.log("user id is " + id);

    return (
        <Container>
            <ToplistCatalog
                title="Toplists by user"
                endpoint={`/toplists/user/${id}`}
            />
        </Container>
    );
}

export default ViewToplistsByUser;
