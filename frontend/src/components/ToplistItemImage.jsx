import PropTypes from "prop-types";
import defaultImage from "../assets/defaultItemImage.jpg";
import { useState, useEffect } from "react";

function ToplistItemImage({ item }) {
    const [loading, setLoading] = useState(true);
    const [src, setSrc] = useState(defaultImage);

    useEffect(() => {
        const getSrc = async () => {
            if (item.newImageFile) {
                setSrc(item.newImageFile);
            } else if (item.image_path) {
                setSrc(
                    `${import.meta.env.VITE_IMG_URL}${
                        item.image_path
                    }?v=${Date.now()}`
                );
            } else {
                setSrc(defaultImage);
            }
        };
        getSrc();
    }, [item]);

    return (
        <div
            style={{
                width: "200px",
                height: "200px",
                overflow: "hidden",
                borderRadius: "5px",
            }}
        >
            <img
                src={src}
                alt={item.title}
                style={{
                    width: "100%",
                    height: "100%",
                    objectFit: "cover",
                    display: loading ? "none" : "block",
                }}
                onLoad={() => {
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
