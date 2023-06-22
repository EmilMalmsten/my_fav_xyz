import './App.css'
import 'bootstrap/dist/css/bootstrap.min.css';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import AnimatedTitle from './components/animated_title';
import ToplistCatalog from './components/toplist_catalog';
import CallToAction from './components/call_to_action';
import MainNavbar from './components/main_navbar';

function App() {

  return (
    <>
      <MainNavbar/>
      <Container>
        <Row className="my-4">
          <Col>
            <AnimatedTitle/>
            <Row>
              <Col>
                <ToplistCatalog title="Most popular toplists"/>
              </Col>
              <Col>
                <ToplistCatalog title="Recent toplists"/>
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

export default App
