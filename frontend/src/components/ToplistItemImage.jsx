import PropTypes from "prop-types";
import defaultImage from "../assets/defaultItemImage.jpg";
import { useState, useEffect } from "react";

function ToplistItemImage({ item }) {
    const [loading, setLoading] = useState(true);
    const [src, setSrc] = useState(defaultImage);

    useEffect(() => {
        const getSrc = async () => {
            let src;
            if (item.newImageFile) {
                src = item.newImageFile;
            } else if (item.image_path) {
                src = `${import.meta.env.VITE_IMG_URL}/${item.list_id}/${
                    item.image_path
                }?v=${Date.now()}`;
            }
            setTimeout(() => {
                setSrc(src);
            }, 400);
        };

        getSrc();
    }, []);

    return (
        <div style={{ width: "200px", height: "200px", overflow: "hidden" }}>
            <img
                src={src}
                alt={item.title}
                style={{
                    width: "100%",
                    height: "100%",
                    objectFit: "cover",
                    display: loading ? "none" : "block",
                }}
                onLoad={(e) => {
                    setLoading(false);
                }}
            />
        </div>
    );
}

ToplistItemImage.propTypes = {
    item: PropTypes.object.isRequired,
};

export default ToplistItemImage;
