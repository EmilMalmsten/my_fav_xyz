import Button from 'react-bootstrap/Button';
import { Link } from 'react-router-dom';

function CallToAction({ title, buttonLink }) {
    return (
      <div className="text-center">
          <h2 className="text-light">{title}</h2>
          <Link to={buttonLink}><Button variant="light" size="lg">Click here</Button></Link>
      </div>
    )
  }
  
  export default CallToAction