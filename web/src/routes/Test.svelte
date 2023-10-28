<script lang="ts">
    import { onMount } from "svelte";
    import { PUBLIC_BASE_URL } from "$env/static/public";

    console.log("Is this the correct URL?: ", PUBLIC_BASE_URL);
</script>

<div>
    URL: {PUBLIC_BASE_URL}
</div>
