import { Row, Col, Container, Button, Card } from "react-bootstrap";
import { useEffect, useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import axios from "axios";

function ToplistSearchResults() {
    const [toplists, setToplists] = useState([]);
    const [searchResultCount, setSearchResultCount] = useState(0);
    const searchParams = new URLSearchParams(location.search);
    const searchTerm = searchParams.get("searchTerm");
    let limit = Number(searchParams.get("limit"));
    let offset = Number(searchParams.get("offset"));
    const navigate = useNavigate();
    const pageSize = 10;

    async function searchToplists() {
        try {
            const resp = await axios.get(
                `${import.meta.env.VITE_API_URL}/toplists/search`,
                {
                    params: {
                        term: searchTerm,
                        limit: limit,
                        offset: offset,
                    },
                }
            );
            setSearchResultCount(resp.data.length);
            setToplists(resp.data.slice(0, pageSize));
        } catch (e) {
            console.error(e);
        }
    }

    function nextPage() {
        limit = 20;
        offset = offset + pageSize;
        navigate(
            `/toplists/search?searchTerm=${searchTerm}&limit=${limit}&offset=${offset}`
        );
        searchToplists();
    }

    function prevPage() {
        limit = 20;
        offset = offset - pageSize;
        navigate(
            `/toplists/search?searchTerm=${searchTerm}&limit=${limit}&offset=${offset}`
        );
        searchToplists();
    }

    useEffect(() => {
        searchToplists();
    }, []);

    return (
        <Container className="my-4">
            <h1>Search results for: {searchTerm}</h1>
            {toplists.length === 0 && <h5>No results found</h5>}
            {toplists.map((toplist) => (
                <Row>
                    <Col>
                        <Link
                            to={`/toplists/${toplist.toplist_id}`}
                            className="link-no-styling"
                        >
                            <h4>{toplist.title}</h4>
                            <p>{toplist.description}</p>
                        </Link>
                    </Col>
                </Row>
            ))}
            {offset > 0 && <Button onClick={prevPage}>Previous page</Button>}
            {searchResultCount > pageSize && (
                <Button onClick={nextPage}>Next page</Button>
            )}
        </Container>
    );
}

export default ToplistSearchResults;
