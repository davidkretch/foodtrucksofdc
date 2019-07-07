import React from "react";

import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faStar } from "@fortawesome/free-solid-svg-icons";

function Star(props) {
    return (
        <FontAwesomeIcon icon={faStar} {...props} />
    );
}

// ButtonStars acts as the rating buttons for the Rating component.
// It also shows light gray placeholder stars.
function ButtonStars(props) {
    const color = "e8e8e8";
    var stars = [];
    for (var i = 1; i <= 3; i++) {
        stars.push(
            <Star
                key={i}
                onMouseOver={props.onMouseOver.bind(this, i)}
                onMouseOut={props.onMouseOut}
                onClick={props.onClick.bind(this, i)}
                style={{color: color, cursor: "pointer"}}
            />
        );
    }
    return (
        <span style={{whiteSpace: "nowrap"}}>
            {stars}
        </span>
    );
}

// DisplayStars shows the rating in stars (either provided or user-selected).
function DisplayStars(props) {
    const color = "ffbc0b";
    const width = Math.floor(100 * props.rating / 3) + "%";
    return (
        <div style={{width: width, overflow: "hidden"}}>
            <span style={{whiteSpace: "nowrap"}}>
                <Star style={{color: color}} />
                <Star style={{color: color}} />
                <Star style={{color: color}} />
            </span>
        </div>
    )
}

class Rating extends React.Component {
    constructor(props) {
        super(props);
        this.highlight = this.highlight.bind(this);
        this.clearHighlight = this.clearHighlight.bind(this);
        this.rate = this.rate.bind(this);
        this.state = {
            rating: this.props.rating || 0,
            selected: null,
            highlighted: null
        }
    }

    highlight(stars) {
        this.setState({
            highlighted: stars
        });
    }

    clearHighlight() {
        this.setState({
            highlighted: null
        });
    }

    rate(stars) {
        this.setState({
            selected: stars
        });
        this.props.setRating(stars);
    }

    render() {
        const rating = this.state.highlighted || this.state.selected || this.state.rating;

        const buttonStars = {
            position: "relative"
        };
        const displayStars = {
            position: "absolute",
            top: 0,
            pointerEvents: "none"
        };

        return (
            <div>
                <div style={buttonStars}>
                    <ButtonStars
                        onMouseOver={this.highlight}
                        onMouseOut={this.clearHighlight}
                        onClick={this.rate}
                    />
                    <div style={displayStars}>
                        <DisplayStars rating={rating} />
                    </div>
                </div>
            </div>
        );
    }
}

export default Rating;