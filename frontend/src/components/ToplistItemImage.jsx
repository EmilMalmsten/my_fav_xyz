import PropTypes from "prop-types";
import defaultImage from "../assets/defaultItemImage.jpg";

function ToplistItemImage({ item }) {
    let src;
    if (item.newImageFile) {
        src = item.newImageFile;
    } else if (item.image_path) {
        src = `${import.meta.env.VITE_IMG_URL}/${item.list_id}/${
            item.image_path
        }?v=${Date.now()}`;
    } else {
        src = defaultImage;
    }

    return (
        <div style={{ width: "200px", height: "200px", overflow: "hidden" }}>
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
