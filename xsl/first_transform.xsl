
<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
    <xsl:output method="html" indent="yes"/>

    <xsl:template match="/">
        <html>
            <head>
                <title>First Transformation</title>
            </head>
            <body>
                <h1>First Transformation Result</h1>
                <xsl:apply-templates/>
            </body>
        </html>
    </xsl:template>

    <xsl:template match="order">
        <div>
            <h2>Order ID: <xsl:value-of select="id"/></h2>
            <p>Description: <xsl:value-of select="description"/></p>
            <p>Amount: <xsl:value-of select="amount"/></p>
        </div>
    </xsl:template>
</xsl:stylesheet>