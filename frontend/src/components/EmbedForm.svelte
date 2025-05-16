{#if data && appliedOverrides}
    <form class="form-wrapper" on:submit|preventDefault>
        <div class="row">
            <Colour col3 label="Embed Farbe" on:change={updateColour} bind:value={tempColour}/>
            <Input col3 label="Titel" placeholder="Embed Title" bind:value={data.title}/>
            <Input col3 label="Titel URL (Optional)" placeholder="https://example.com" bind:value={data.url}/>
        </div>

        <div class="row">
            <Textarea col1 label="Beschreibung" placeholder="Großer Text Bereich, bis zu 4096 Zeichen"
                      bind:value={data.description}/>
        </div>

        <Collapsible forceAlwaysOpen>
            <span slot="header">Author</span>

            <div slot="content" class="row">
                <Input col3 label="Author Name" placeholder="Author Name" bind:value={data.author.name}/>
                <Input col3 label="Author Icon URL (Optional)" placeholder="https://example.com/image.png"
                       tooltipText="Kleines Icon was oben links angezeigt wird" bind:value={data.author.icon_url}/>
                <Input col3 label="Author URL (Optional)" placeholder="https://example.com"
                       tooltipText="Link auf dem Authorennamen" bind:value={data.author.url}/>
            </div>
        </Collapsible>

        <Collapsible forceAlwaysOpen>
            <span slot="header">Bilder</span>
            <div slot="content" class="row">
                <Input col2 label="Großes Bild URL" placeholder="https://example.com/image.png"
                       bind:value={data.image_url}/>
                <Input col2 label="Kleines Bild URL" placeholder="https://example.com/image.png"
                       bind:value={data.thumbnail_url}/>
            </div>
        </Collapsible>

        <Collapsible forceAlwaysOpen>
            <span slot="header">Fußzeile</span>
            <div slot="content" class="row">
                {#if footerPremiumOnly}
                    <Input col3 label="Fußzeilen Text" placeholder="Fußzeilen Text" badge="Premium"
                           bind:value={data.footer.text}/>
                    <Input col3 label="Fußzeilen Icon URL (Optional)" badge="Premium"
                           placeholder="https://example.com/image.png"
                           bind:value={data.footer.icon_url}/>
                {:else}
                    <Input col3 label="Fußzeilen Text" placeholder="Footer Text" bind:value={data.footer.text}/>
                    <Input col3 label="Fußzeilen Icon URL (Optional)" placeholder="https://example.com/image.png"
                           bind:value={data.footer.icon_url}/>
                {/if}
                <DateTimePicker col3 label="Fußzeilen Timestamp (Optional)" bind:value={data.timestamp}/>
            </div>
        </Collapsible>

        <Collapsible forceAlwaysOpen>
            <span slot="header">Felder</span>
            <div slot="content" class="col-1">
                {#each data.fields as field, i}
                    <div class="row" style="justify-content: flex-start; gap: 10px">
                        <Input col2 label="Feld Name" placeholder="Field Name" bind:value={field.name}/>
                        <Checkbox label="Einzeilig" bind:value={field.inline}/>

                        <div style="margin-top: 18px; display: flex; align-self: center">
                            <Button danger icon="fas fa-trash-can" on:click={() => deleteField(i)}>Löschen</Button>
                        </div>
                    </div>
                    <div class="row">
                      <Textarea col1 label="Feld Wert" placeholder="Großer Text bereich, bis zu 1024 Zeichen"
                                bind:value={field.value}/>
                    </div>
                {/each}

                <div class="add-field-wrapper">
                    <Button type="button" icon="fas fa-plus" fullWidth on:click={addField}>Feld hinzufügen</Button>
                </div>
            </div>
        </Collapsible>
    </form>
{/if}

<style>
    .form-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
    }

    .row {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        gap: 10px;
        width: 100%;
    }

    .add-field-wrapper {
        width: 100%;
        margin-top: 10px;
    }
</style>

<script>
    import Textarea from "./form/Textarea.svelte";
    import Colour from "./form/Colour.svelte";
    import Input from "./form/Input.svelte";
    import Collapsible from "./Collapsible.svelte";
    import DateTimePicker from "./form/DateTimePicker.svelte";
    import Checkbox from "./form/Checkbox.svelte";
    import Button from "./Button.svelte";
    import {onMount} from "svelte";
    import {intToColour, colourToInt} from "../js/util";

    export let data;

    $: data = data ?? {
        fields: [],
        colour: 0x2ECC71,
        author: {},
        footer: {},
    };

    export let footerPremiumOnly = true;

    function addField() {
        data.fields.push({name: '', value: '', inline: false});
        data = data;
    }

    function deleteField(i) {
        data.fields.splice(i, 1);
        data = data;
    }

    let tempColour = "#2ecc71";
    function updateColour() {
        data.colour = colourToInt(tempColour);
    }

    let appliedOverrides = false;
    onMount(() => {
        data.author = data.author ?? {};
        data.footer = data.footer ?? {};
        data.fields = data.fields ?? [];

        if (!data.colour) {
            data.colour = 0x2ECC71;
        } else {
            if (typeof data.colour === "string" && data.colour.startsWith('#')) {
                data.colour = parseInt(data.colour.slice(1), 16)
            }

            tempColour = intToColour(data.colour);
        }

        appliedOverrides = true;
    });
</script>
