import AnimatedTitle from "../components/AnimatedTitle";
import ToplistCatalog from "../components/ToplistCatalog";
import { Button, Container, Alert } from "react-bootstrap";
import { useLocation } from "react-router-dom";

function Home() {
    const location = useLocation();
    const successAlert = location.state && location.state.successAlert;

    return (
        <Container style={{ width: "80%", margin: "0 auto" }}>
            <div className="my-3">
                {successAlert && (
                    <Alert variant="success" dismissible>
                        {successAlert}
                    </Alert>
                )}
            </div>
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
