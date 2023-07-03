import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import AnimatedTitle from "../components/AnimatedTitle";
import ToplistCatalog from "../components/ToplistCatalog";
import CallToAction from "../components/CallToAction";

function Home() {
    return (
        <>
            <Container>
                <Row className="my-4">
                    <Col>
                        <AnimatedTitle />
                        <Row>
                            <Col className="my-4">
                                <ToplistCatalog
                                    title="Most popular toplists"
                                    endpoint="/toplists/popular"
                                />
                            </Col>
                            <Col className="my-4">
                                <ToplistCatalog
                                    title="Recent toplists"
                                    endpoint="/toplists/recent"
                                />
                            </Col>
                        </Row>
                    </Col>
                    <Col className="bg-secondary d-flex align-items-center justify-content-center">
                        <CallToAction
                            title="Create your own toplist"
                            buttonLink="/toplists/create"
                        />
                    </Col>
                </Row>
            </Container>
        </>
    );
}

export default Home;
