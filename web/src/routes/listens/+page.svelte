<script lang="ts">
</script>

<div class="grid">
    <h1>Music</h1>
    <div class="foo">
        <ul class="masonry">
            <li class="masonry-brick">foo</li>
            <li class="masonry-brick">foo</li>
            <li class="masonry-brick">foo</li>
            <li class="masonry-brick">foo</li>
            <li class="masonry-brick">foo</li>
            <li class="masonry-brick">foo</li>
        </ul>
    </div>
</div>

<style>
    .grid {
        display: grid;
        grid-template-rows: 100px 1fr;
        height: 100%;
        width: 100%;
    }

    .masonry {
        display: grid;
        grid-auto-flow: dense;
        grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
        justify-items: center;
        align-items: center;
        background-size: cover;
        object-fit: cover;
        gap: 10px;
    }

    .masonry > li {
        list-style: none;
        background-color: bisque;
        height: 200px;
        width: 200px;
        margin: 1;
    }

    .masonry-brick {
        object-fit: cover;
        height: inherit;
        max-width: 100%;
    }
</style>
