<script>
    import profile from "$lib/images/pfp.webp";
</script>

<section>
    <div class="left">
        <img src={profile} alt="Welcome" style="width: 250px;" />
        <h1>Hi, I'm Josh</h1>
    </div>
    <div class="right">
        <article>
            <h2>Welcome to my portfolio.</h2>
    
            <p>
                I'm primarily a full stack developer, but my interests span
                everything from embedded sytems to dev-ops.
            </p>
            <p>
                Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Condimentum lacinia quis vel eros donec. Egestas fringilla phasellus faucibus scelerisque eleifend donec pretium vulputate sapien. Massa tincidunt dui ut ornare lectus sit amet. Rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt lobortis feugiat. Malesuada proin libero nunc consequat interdum varius sit. Ornare suspendisse sed nisi lacus sed. Aliquam faucibus purus in massa tempor nec. Phasellus faucibus scelerisque eleifend donec pretium. At elementum eu facilisis sed odio morbi. Praesent elementum facilisis leo vel fringilla est.
                Duis at consectetur lorem donec massa sapien faucibus et. Tortor condimentum lacinia quis vel eros donec. Ultrices gravida dictum fusce ut placerat orci nulla pellentesque dignissim. Sed vulputate mi sit amet. Leo vel orci porta non pulvinar neque. Eu ultrices vitae auctor eu augue ut lectus arcu bibendum. Purus semper eget duis at tellus. Augue neque gravida in fermentum et sollicitudin ac orci. Pharetra sit amet aliquam id diam maecenas ultricies mi eget. Ac tortor vitae purus faucibus ornare. Interdum velit laoreet id donec. Dictum non consectetur a erat nam at lectus. Libero justo laoreet sit amet cursus sit amet dictum sit. Amet volutpat consequat mauris nunc congue nisi vitae suscipit tellus.
            </p>
        </article>

        
        <article>
            <h3>Languages</h3>
            <ul>
                <li>c++</li>
                <li>javascript</li>
                <li>typescript</li>
                <li>python</li>

            </ul>

            <h5>Learning</h5>
            <ul>
                <li>rust</li>
                <li>go</li>
                <li>shell</li>
            </ul>

            <h3>Frameworks</h3>
            <ul>
                <li>vue</li>
                <li>svelte</li>
                <li>node</li>
                <li>express</li>
                <li>InfluxDB</li>
                <li>mysql</li>
            </ul>

            <h3>Technologies</h3>
            <ul>
                <li>linux</li>
                <li>dev-ops</li>
                <li>docker</li>
                <li>CI/CD</li>
                <li>git</li>
                <li>security</li>
            </ul>
        </article>
    </div>


</section>

<style>
    h1, h2, h3, h5 {
        margin: 0px;
    }

    section {
        display: flex;
        gap: 10px;
    }

    section > .left {
        width: 250px;
        font-size: 5em;
        text-align: justify;
    }

    section > .right {
        display: flex;
        flex-direction: column;
        gap: 20px;
    }

    article {
        font-size: 0.9em;
        color: white;
        width: 400px;
        height: 50vh;
        overflow: auto;
        text-align: justify;
    }

    ul {
        display: flex;
        flex-direction: row;
        gap: 10px;
        list-style: none;
        margin: 0px;
    }
</style>
