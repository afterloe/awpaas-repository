"use strict";

/**
 *  顶部导航栏
 */
class Header extends React.Component {
    constructor(props) {
        super(props);
        this.blurActive = this.blurActive.bind(this);
        this.getFocus = this.getFocus.bind(this);
    }

    blurActive(event) {
        const dom = event.currentTarget;
        dom.setAttribute("class", "input")
    }

    getFocus(event) {
        const dom = event.currentTarget;
        dom.setAttribute("class", "input isActive")
    }

    render() {
        let p = this.props;
        return (
            <nav class="navbar navbar-dark fixed-top bg-dark flex-md-nowrap p-0 shadow">
                <a class="navbar-brand col-sm-3 col-md-2 mr-0" href="#">{p.name}</a>
                <input class="form-control form-control-dark w-100 input" type="text" placeholder="搜索" aria-label="搜索"
                       onFocus={this.getFocus} onBlur={this.blurActive}/>
                <ul class="navbar-nav px-3">
                    <li class="nav-item text-nowrap">
                        <a class="nav-link" href="#">{p.version}</a>
                    </li>
                </ul>
            </nav>
    )
    }
}

ReactDOM.render(<Header name="镜像管理" version="v1.0.0"/>, document.getElementById("head"));