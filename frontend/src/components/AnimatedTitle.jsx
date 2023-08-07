import React, { useState, useEffect } from "react";

function AnimatedTitle() {
    const texts = [
        "Restaurants in London",
        "Cat cafes in Tokyo",
        "10 favorite movies",
        "Baking recipes",
    ];
    const [index, setIndex] = useState(0);
    const [subIndex, setSubIndex] = useState(0);
    const [isErasing, setIsErasing] = useState(false);

    useEffect(() => {
        if (isErasing && subIndex > 0) {
            setTimeout(() => {
                setSubIndex(subIndex - 1);
            }, 50);
        } else if (isErasing && subIndex === 0) {
            setIsErasing(false);
            if (index < texts.length - 1) {
                setIndex(index + 1);
            } else {
                setIndex(0);
            }
        } else if (subIndex < texts[index].length + 1) {
            setTimeout(() => {
                setSubIndex(subIndex + 1);
            }, 150);
        } else if (subIndex === texts[index].length + 1) {
            setTimeout(() => {
                setIsErasing(true);
            }, 1000);
        }
    }, [subIndex, index, isErasing, texts]);

    return (
        <div className="my-4 text-center">
            <div className="d-flex align-items-center justify-content-center">
                <h2 className="mb-0">My Top </h2>
                <span className="dynamic-text">
                    {texts[index].substring(0, subIndex)}
                    <span className={"cursor blink"}>|</span>
                </span>
            </div>
        </div>
    );
}

export default AnimatedTitle;
