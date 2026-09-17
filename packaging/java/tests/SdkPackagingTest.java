package com.htmlcsstoimage.pulumi;

import com.google.gson.JsonParser;
import org.junit.jupiter.api.Test;
import java.io.InputStreamReader;
import static org.junit.jupiter.api.Assertions.*;

class SdkPackagingTest {
    @Test
    void runtimeFindsVersionAndPluginMetadata() {
        assertEquals(System.getProperty("sdk.version"), Utilities.getVersion());
        var stream = Utilities.class.getClassLoader().getResourceAsStream(
                "com/htmlcsstoimage/html-css-to-image/plugin.json");
        assertNotNull(stream, "Plugin discovery metadata must ship in the SDK");
        var metadata = JsonParser.parseReader(new InputStreamReader(stream)).getAsJsonObject();
        assertEquals(Utilities.getVersion(), metadata.get("version").getAsString());
        assertEquals("html-css-to-image", metadata.get("name").getAsString());
        assertEquals("github://api.github.com/htmlcsstoimage/pulumi-html-css-to-image",
                metadata.get("server").getAsString());
        assertTrue(metadata.get("resource").getAsBoolean());
    }
}
