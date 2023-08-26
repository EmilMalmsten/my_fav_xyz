import { useLocation } from "react-router-dom";
import { Row, Col, Container } from "react-bootstrap";

function ListToplists() {
    const location = useLocation();
    const toplists = location.state || {};
    const searchParams = new URLSearchParams(location.search);
    const searchTerm = searchParams.get("searchTerm");
    const page = searchParams.get("page");

    return (
        <Container>
            <h1>Search results for: {searchTerm}</h1>
            {toplists.map((toplist) => (
                <Row>
                    <Col>
                        <h4>{toplist.title}</h4>
                        <p>{toplist.description}</p>
                    </Col>
                </Row>
            ))}
        </Container>
    );
}

export default ListToplists;
