import Button from "react-bootstrap/Button";
import { Link } from "react-router-dom";
import PropTypes from "prop-types";

function CallToAction({ title, buttonLink }) {
    return (
        <div className="text-center">
            <h2 className="text-light">{title}</h2>
            <Link to={buttonLink}>
                <Button variant="light" size="lg">
                    Click here
                </Button>
            </Link>
        </div>
    );
}

CallToAction.propTypes = {
    title: PropTypes.string.isRequired,
    buttonLink: PropTypes.string.isRequired,
};

export default CallToAction;
