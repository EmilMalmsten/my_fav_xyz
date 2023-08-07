import AnimatedTitle from "../components/AnimatedTitle";
import ToplistCatalog from "../components/ToplistCatalog";
import { Button, Container, Col, Row } from "react-bootstrap";

function Home() {
    return (
        <Container style={{ width: "80%", margin: "0 auto" }}>
            <AnimatedTitle />

            <ToplistCatalog
                title="Most popular toplists"
                endpoint="/toplists/popular"
            />

            <ToplistCatalog
                title="Recent toplists"
                endpoint="/toplists/recent"
            />

            <div
                style={{ textAlign: "center", backgroundColor: "lightGray" }}
                className="my-3 py-3"
            >
                <h5>Make your own toplist!</h5>
                <Button>Create now</Button>
            </div>
        </Container>
    );
}

export default Home;
