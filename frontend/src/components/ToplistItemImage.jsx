import PropTypes from "prop-types";
import defaultImage from "../assets/defaultItemImage.jpg";

function ToplistItemImage({ item }) {
    console.log(item);
    let src;
    if (item.newImageURL) {
        src = item.newImageURL;
    } else if (item.image_path) {
        src = `http://localhost:8080/images/${item.list_id}/${item.image_path}`;
    } else {
        src = defaultImage;
    }

    return (
        <div style={{ width: "100px", height: "100px", overflow: "hidden" }}>
            <img
                src={src}
                alt={item.title}
                style={{ width: "100%", height: "100%", objectFit: "cover" }}
            />
        </div>
    );
}

ToplistItemImage.propTypes = {
    item: PropTypes.object.isRequired,
};

export default ToplistItemImage;
