import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import AnimatedTitle from './components/animated_title';
import ToplistCatalog from './components/toplist_catalog';
import CallToAction from './components/call_to_action';

function Home() {
  return (
    <>
      <Container>
        <Row className="my-4">
          <Col>
            <AnimatedTitle/>
            <Row>
              <Col className="my-4">
                <ToplistCatalog title="Most popular toplists" endpoint="/toplists"/>
              </Col>
              <Col className="my-4">
                <ToplistCatalog title="Recent toplists" endpoint="/toplists/recent"/>
              </Col>
            </Row>
          </Col>
          <Col className="bg-secondary d-flex align-items-center justify-content-center">
            <CallToAction/>
          </Col>
        </Row>
      </Container>
    </>
  )
}

export default Home